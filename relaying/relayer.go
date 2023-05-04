package relaying

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/relaying/log"
)

type IRelayer interface {
	SendRelayTransaction(args RelayTransactionArgs) (common.Hash, error)
	GetRelayReceipt(txHash common.Hash) (RelayReceipt, error)
	GetRefundTokens() []common.Address
	GetRefundAddresses() []common.Address
}

var _ IRelayer = (*Relayer)(nil)

type Relayer struct {
	config         RelayerConfig
	logger         log.Logger
	abi            *abi.ABI
	relayerManager *RelayerManager
	provider       IProvider
	keyStore       IKeyStore

	pendingSenders  map[string]struct{}
	nextSenderIndex int
	mtx             *sync.Mutex
}

func NewRelayer(
	logger log.Logger,
	config RelayerConfig,
	client bind.ContractBackend,
	provider IProvider,
	keyStore IKeyStore,
) (IRelayer, error) {
	parsedABI, err := RelayerManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	relayerManager, err := NewRelayerManager(config.Address, client)
	if err != nil {
		return nil, err
	}

	return &Relayer{
		config:          config,
		logger:          logger,
		abi:             parsedABI,
		relayerManager:  relayerManager,
		provider:        provider,
		keyStore:        keyStore,
		pendingSenders:  make(map[string]struct{}),
		nextSenderIndex: 0,
		mtx:             &sync.Mutex{},
	}, nil
}

func NewRelayCall(args RelayTransactionArgs) RelayTransaction {
	return RelayTransaction{
		Wallet:        ensureAddress(args.Wallet),
		Data:          ensureBytes(args.Data),
		Nonce:         (*big.Int)(args.Nonce),
		Signatures:    ensureBytes(args.Signatures),
		GasPrice:      (*big.Int)(args.GasPrice),
		GasLimit:      (*big.Int)(args.GasLimit),
		RefundToken:   ensureAddress(args.RefundToken),
		RefundAddress: ensureAddress(args.RefundAddress),
	}
}

func (relayer *Relayer) SendRelayTransaction(args RelayTransactionArgs) (common.Hash, error) {
	call := NewRelayCall(args)

	err := relayer.validateCall(call)
	if err != nil {
		return common.Hash{}, NewRpcErrorf(ErrorCodeInvalidRequest, nil, "invalid request: %s", err)
	}

	senderAddress, err := relayer.getNextSenderRetry(5, time.Second)
	if err != nil {
		return common.Hash{}, NewRpcErrorf(ErrorCodeSystemBusy, nil, "could not get next sender address")
	}
	relayer.setSenderStatus(senderAddress, true)
	defer relayer.setSenderStatus(senderAddress, false)

	callData, err := relayer.abi.Pack("execute", call.Wallet,
		call.Data,
		call.Nonce,
		call.Signatures,
		call.GasPrice,
		call.GasLimit,
		call.RefundToken,
		call.RefundAddress,
	)

	estimateGas, err := relayer.provider.EstimateGas(
		ethereum.CallMsg{
			From: senderAddress,
			To:   &relayer.config.Address,
			Data: callData,
		},
	)
	if err != nil {
		return common.Hash{}, NewRpcErrorf(ErrorCodeInvalidRequest, nil, "could not estimate gas: %s", err)
	}

	gasLimit := uint64(float64(estimateGas) * relayer.config.GasMultiplier)

	tx, err := relayer.relayerManager.Execute(
		&bind.TransactOpts{
			NoSend: true,
			From:   senderAddress,
			Signer: func(_ common.Address, transaction *ethtypes.Transaction) (*ethtypes.Transaction, error) {
				return relayer.keyStore.SignTx(senderAddress, transaction)
			},
			GasLimit: gasLimit,
		},
		call.Wallet,
		call.Data,
		call.Nonce,
		call.Signatures,
		call.GasPrice,
		call.GasLimit,
		call.RefundToken,
		call.RefundAddress,
	)
	if err != nil {
		return common.Hash{}, NewRpcErrorf(ErrorCodeInternal, nil, "could not build tx: %s", err)
	}

	rawTx, _ := tx.MarshalBinary()
	err = relayer.provider.SendRawTransaction(rawTx)
	if err != nil {
		return common.Hash{}, NewRpcErrorf(ErrorCodeInternal, nil, "could not send relay transaction: %s", err)
	}

	return tx.Hash(), nil
}

func (relayer *Relayer) getRefundTokenIndex(tokenAddr common.Address) int {
	for idx, addr := range relayer.config.AllowedRefundTokens {
		if addr == tokenAddr {
			return idx
		}
	}
	return -1
}

func (relayer *Relayer) validateCall(call RelayTransaction) error {
	if call.Nonce == nil {
		return errors.New("missing nonce field")
	}
	if call.GasPrice == nil {
		return errors.New("missing gasPrice field")
	}
	if call.GasLimit == nil {
		return errors.New("missing gasLimit field")
	}

	if !addressInArray(relayer.config.AllowedRefundAddresses, call.RefundAddress) {
		return fmt.Errorf("refund address not allowed")
	}

	refundTokenIdx := relayer.getRefundTokenIndex(call.RefundToken)
	if refundTokenIdx == -1 {
		return fmt.Errorf("refund token not allowed")
	}
	if call.GasPrice.Cmp(relayer.config.MinGasPrices[refundTokenIdx]) < 0 {
		return fmt.Errorf("gasPrice too low, expected at least %s", relayer.config.MinGasPrices[refundTokenIdx].String())
	}

	return nil
}

var transactionExecutedEventSignature = crypto.Keccak256Hash([]byte(`TransactionExecuted(address,bool,bytes,bytes32)`))

type RelayReceipt struct {
	TransactionHash common.Hash `json:"transactionHash"`
	SignHash        common.Hash `json:"signHash"`
	Success         bool        `json:"success"`
	Error           string      `json:"error,omitempty"`
}

func (relayer *Relayer) GetRelayReceipt(txHash common.Hash) (RelayReceipt, error) {
	receipt, err := relayer.provider.GetTransactionReceipt(txHash)
	if err != nil {
		return RelayReceipt{}, NewRpcErrorf(ErrorCodeInvalidRelayTx, nil, "could not get transaction receipt: %s", err)
	}

	executeLogIdx := -1
	for idx, ethLog := range receipt.Logs {
		if ethLog == nil || len(ethLog.Topics) == 0 {
			continue
		}
		if ethLog.Topics[0] == transactionExecutedEventSignature {
			executeLogIdx = idx
			break
		}
	}

	if executeLogIdx == -1 {
		return RelayReceipt{}, NewRpcErrorf(ErrorCodeInvalidRelayTx, nil, "could not find relay event in transaction receipt")
	}

	evt, err := relayer.relayerManager.ParseTransactionExecuted(*receipt.Logs[executeLogIdx])
	if err != nil {
		return RelayReceipt{}, NewRpcErrorf(ErrorCodeInvalidRelayTx, nil, "could not parse event log: %s", err)
	}

	var relayErr string
	if !evt.Success && len(evt.ReturnData) > 0 {
		var errUnpack error
		relayErr, errUnpack = abi.UnpackRevert(evt.ReturnData)
		if errUnpack != nil {
			relayErr = common.Bytes2Hex(evt.ReturnData)
		}
	}

	return RelayReceipt{
		TransactionHash: txHash,
		SignHash:        evt.SignedHash,
		Success:         evt.Success,
		Error:           relayErr,
	}, nil
}

func (relayer *Relayer) WaitForTransaction(ctx context.Context, txHash common.Hash) (*ethtypes.Receipt, error) {
	queryTicker := time.NewTicker(500 * time.Millisecond)
	defer queryTicker.Stop()

	for {
		receipt, err := relayer.provider.GetTransactionReceipt(txHash)
		if err == nil {
			return receipt, nil
		}

		if errors.Is(err, ethereum.NotFound) {
			relayer.logger.Debug("transaction not yet mined")
		} else {
			relayer.logger.Debug("receipt retrieval failed", "err", err)
		}

		// Wait for the next round.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}

func (relayer *Relayer) GetRefundAddresses() []common.Address {
	return relayer.config.AllowedRefundAddresses
}

func (relayer *Relayer) GetRefundTokens() []common.Address {
	return relayer.config.AllowedRefundTokens
}

var errSendersBusy = errors.New("senders are busy")

func (relayer *Relayer) getNextSender() (common.Address, error) {
	relayer.mtx.Lock()
	defer relayer.mtx.Unlock()

	idx := relayer.nextSenderIndex
	l := len(relayer.config.SenderAddresses)

	// return immediately if the next sender is available
	if _, ok := relayer.pendingSenders[relayer.config.SenderAddresses[idx].Hex()]; !ok {
		relayer.nextSenderIndex = (idx + 1) % l
		return relayer.config.SenderAddresses[idx], nil
	}

	// find next sender that is available
	for {
		if idx == relayer.nextSenderIndex {
			return common.Address{}, errSendersBusy
		}
		if _, ok := relayer.pendingSenders[relayer.config.SenderAddresses[idx].Hex()]; !ok {
			break
		}
		idx = (idx + 1) % l
	}

	relayer.nextSenderIndex = (idx + 1) % l
	return relayer.config.SenderAddresses[idx], nil
}

func (relayer *Relayer) getNextSenderRetry(maxRetry int, delay time.Duration) (common.Address, error) {
	for attempt := 1; attempt <= maxRetry; attempt++ {
		addr, err := relayer.getNextSender()
		if err == nil {
			return addr, nil
		}
		time.Sleep(delay)
	}
	return common.Address{}, errSendersBusy
}

func (relayer *Relayer) setSenderStatus(sender common.Address, busy bool) {
	relayer.mtx.Lock()
	defer relayer.mtx.Unlock()

	if busy {
		relayer.pendingSenders[sender.Hex()] = struct{}{}
	} else {
		delete(relayer.pendingSenders, sender.Hex())
	}
}
