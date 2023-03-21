package eip4337

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	entrypoint_interface "github.com/evmos/ethermint/eip4337/entrypoint"
	"github.com/evmos/ethermint/eip4337/log"
	"github.com/evmos/ethermint/eip4337/types"
	"github.com/pkg/errors"
)

type BundleManager struct {
	logger                 log.Logger
	beneficiary            common.Address
	minSignerBalance       *big.Int
	maxBundleGas           uint64
	signer                 ISigner
	provider               IProvider
	eventsManager          IEventsManager
	mempoolManager         IMempoolManager
	reputationManager      IReputationManager
	validationManager      IValidationManager
	entryPoint             entrypoint_interface.IEntryPoint
	mergeToAccountRootHash bool
	conditionalRpc         bool
	mtx                    *sync.RWMutex
}

func NewBundleManager(
	logger log.Logger,
	provider IProvider,
	signer ISigner,
	entryPoint entrypoint_interface.IEntryPoint,
	eventsManager IEventsManager,
	mempoolManager IMempoolManager,
	validationManager IValidationManager,
	reputationManager IReputationManager,
	beneficiary common.Address,
	minSignerBalance *big.Int,
	maxBundleGas uint64,
) *BundleManager {
	return &BundleManager{
		logger:                 log.EnsureLogger(logger),
		beneficiary:            beneficiary,
		minSignerBalance:       minSignerBalance,
		maxBundleGas:           maxBundleGas,
		signer:                 signer,
		provider:               provider,
		eventsManager:          eventsManager,
		mempoolManager:         mempoolManager,
		reputationManager:      reputationManager,
		validationManager:      validationManager,
		entryPoint:             entryPoint,
		mergeToAccountRootHash: false,
		conditionalRpc:         false,
		mtx:                    &sync.RWMutex{},
	}
}

var ErrNotEnoughOps = errors.New("not enough user operations to bundle")

type SendBundleReturn struct {
	TransactionHash common.Hash
	UserOpHashes    []common.Hash
}

func (manager *BundleManager) SendNextBundle() (SendBundleReturn, error) {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	err := manager.HandlePastEvents()
	if err != nil {
		return SendBundleReturn{}, err
	}

	bundle := manager.createBundle()
	if len(bundle) == 0 {
		return SendBundleReturn{}, ErrNotEnoughOps
	}

	beneficiary, err := manager.selectBeneficiary()
	if err != nil {
		return SendBundleReturn{}, err
	}

	return manager.sendBundle(bundle, beneficiary)
}

func ToABIUserOperation(op types.UserOperation) entrypoint_interface.UserOperation {
	return entrypoint_interface.UserOperation{
		Sender:               *op.Sender(),
		Nonce:                op.Nonce(),
		InitCode:             op.InitCode(),
		CallData:             op.CallData(),
		CallGasLimit:         op.CallGasLimit(),
		VerificationGasLimit: op.VerificationGasLimit(),
		PreVerificationGas:   op.PreVerificationGas(),
		MaxFeePerGas:         op.MaxFeePerGas(),
		MaxPriorityFeePerGas: op.MaxPriorityFeePerGas(),
		PaymasterAndData:     op.PaymasterAndData(),
		Signature:            op.Signature(),
	}
}

func ToABIUserOperations(ops []types.UserOperation) []entrypoint_interface.UserOperation {
	var ret []entrypoint_interface.UserOperation
	for _, op := range ops {
		ret = append(ret, ToABIUserOperation(op))
	}
	return ret
}

func (manager *BundleManager) sendBundle(bundle []types.UserOperation, beneficiary common.Address) (SendBundleReturn, error) {
	bundle, tx, err := manager.prepareBundleTx(bundle, beneficiary)
	if err != nil {
		return SendBundleReturn{}, err
	}

	rawTx, _ := tx.MarshalBinary()
	err = manager.provider.SendRawTransaction(rawTx)
	if err != nil {
		return SendBundleReturn{}, err
	}

	// use context.WithTimeout to add deadline
	receipt, err := manager.waitMined(context.Background(), tx)
	if err != nil {
		return SendBundleReturn{}, err
	}

	if receipt.Status != ethtypes.ReceiptStatusSuccessful {
		// ignore error
		_, _ = manager.handleFailedOps(bundle, beneficiary)
		return SendBundleReturn{}, fmt.Errorf("bundle tx failed. txHash = %s", tx.Hash())
	}

	// remove ops from mempool
	for _, op := range bundle {
		manager.mempoolManager.RemoveUserOp(op)
	}

	chainID, err := manager.provider.ChainID()
	if err != nil {
		return SendBundleReturn{}, err
	}

	hashes := GetUserOpHashes(bundle, manager.entryPoint.Address(), chainID)

	return SendBundleReturn{
		TransactionHash: tx.Hash(),
		UserOpHashes:    hashes,
	}, nil
}

func (manager *BundleManager) prepareBundleTx(bundle []types.UserOperation, beneficiary common.Address) ([]types.UserOperation, *ethtypes.Transaction, error) {
	for {
		gasPrice, err := manager.provider.GetGasPrice()
		if err != nil {
			return bundle, nil, err
		}

		tx, err := manager.entryPoint.Transactor().HandleOps(&bind.TransactOpts{
			NoSend:   true,
			GasPrice: gasPrice,
			Signer: func(_ common.Address, transaction *ethtypes.Transaction) (*ethtypes.Transaction, error) {
				return manager.signer.SignTx(transaction)
			},
		}, ToABIUserOperations(bundle), beneficiary)
		if err == nil {
			return bundle, tx, nil
		}

		// a workaround for cosmos-sdk estimateGas does not return what we need
		if !strings.Contains(err.Error(), "execution revert") {
			return bundle, nil, err
		}
		bundle, err = manager.handleFailedOps(bundle, beneficiary)
		if err != nil {
			return bundle, nil, errors.Wrapf(err, "estimate gas failed, but cannot handle the error")
		}

		if len(bundle) == 0 {
			return bundle, nil, ErrNotEnoughOps
		}
	}
}

func (manager *BundleManager) handleFailedOps(bundle []types.UserOperation, beneficiary common.Address) ([]types.UserOperation, error) {
	err := manager.entryPoint.Caller().HandleOps(&bind.CallOpts{}, ToABIUserOperations(bundle), beneficiary)
	if err == nil {
		return bundle, errors.New("call not failed")
	}
	if strings.Contains(err.Error(), "execution reverted") {
		return bundle, err
	}
	failedOpError, decodeErr := manager.entryPoint.ErrorDecoder().DecodeFailedOp(err)
	if err != nil {
		return bundle, decodeErr
	}
	if failedOpError.OpIndex.Cmp(new(big.Int).SetInt64(int64(len(bundle)))) >= 0 {
		// should never happen
		return bundle, fmt.Errorf("invalid opIndex returned %d", failedOpError.OpIndex.Int64())
	}

	failedUserOp := bundle[failedOpError.OpIndex.Uint64()]
	// Remove the failed op that caused the revert from the batch and drop from the mempool.
	// Other ops from the same paymaster should be removed from the current batch, but kept in the mempool
	manager.mempoolManager.RemoveUserOp(failedUserOp)
	bundle = removeFailedUserOpByIndex(bundle, failedOpError.OpIndex.Uint64())
	switch {
	case strings.HasPrefix(failedOpError.Reason, "AA3"):
		manager.reputationManager.CrashedHandleOps(getAddr(failedUserOp.PaymasterAndData()))
	case strings.HasPrefix(failedOpError.Reason, "AA2"):
		manager.reputationManager.CrashedHandleOps(failedUserOp.Sender())
	case strings.HasPrefix(failedOpError.Reason, "AA1"):
		manager.reputationManager.CrashedHandleOps(getAddr(failedUserOp.InitCode()))
	default:
		return bundle, fmt.Errorf("unknown error: %s", failedOpError.Reason)
	}

	return bundle, nil
}

func removeFailedUserOpByIndex(userOps []types.UserOperation, idx uint64) []types.UserOperation {
	failedOp := userOps[idx]
	failedPaymaster := getAddr(failedOp.PaymasterAndData())
	if failedPaymaster == nil {
		return append(userOps[:idx], userOps[idx+1:]...)
	}

	newOps := make([]types.UserOperation, 0)
	for _, op := range userOps {
		paymaster := getAddr(op.PaymasterAndData())
		if paymaster != nil && paymaster.String() == failedPaymaster.String() {
			continue
		}
		newOps = append(newOps, op)
	}

	return newOps
}

func (manager *BundleManager) waitMined(ctx context.Context, tx *ethtypes.Transaction) (*ethtypes.Receipt, error) {
	queryTicker := time.NewTicker(500 * time.Millisecond)
	defer queryTicker.Stop()

	for {
		receipt, err := manager.provider.GetTransactionReceipt(tx.Hash())
		if err == nil {
			return receipt, nil
		}

		if errors.Is(err, ethereum.NotFound) {
			manager.logger.Debug("Transaction not yet mined")
		} else {
			manager.logger.Debug("Receipt retrieval failed", "err", err)
		}

		// Wait for the next round.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}

func (manager *BundleManager) selectBeneficiary() (common.Address, error) {
	currentBalance, err := manager.provider.GetBalance(manager.signer.GetAddress())
	if err != nil {
		return common.Address{}, err
	}
	beneficiary := manager.beneficiary
	if currentBalance.Cmp(manager.minSignerBalance) <= 0 {
		beneficiary = manager.signer.GetAddress()
	}
	return beneficiary, nil
}

func (manager *BundleManager) HandlePastEvents() error {
	return manager.eventsManager.HandlePastEvents()
}

func (manager *BundleManager) createBundle() []types.UserOperation {
	entries := manager.mempoolManager.GetSortedForInclusion()
	bundle := make([]types.UserOperation, 0)

	paymasterDeposit := make(map[string]*big.Int)
	stakedEntityCount := make(map[string]int64)
	senders := make(map[string]struct{})
	totalGas := big.NewInt(0)

	for _, entry := range entries {
		paymaster := getAddr(entry.UserOp.PaymasterAndData())
		factory := getAddr(entry.UserOp.InitCode())
		paymasterStatus := manager.reputationManager.Status(paymaster)
		deployerStatus := manager.reputationManager.Status(factory)
		if paymasterStatus == ReputationStatusBanned || deployerStatus == ReputationStatusBanned {
			manager.mempoolManager.RemoveUserOp(entry.UserOp)
		}
		if paymaster != nil && paymasterStatus == ReputationStatusThrottled && stakedEntityCount[paymaster.String()] > 1 {
			manager.logger.Debug("skipping throttled paymaster", "sender", entry.UserOp.Sender(), "nonce", entry.UserOp.Nonce())
			continue
		}
		if factory != nil && deployerStatus == ReputationStatusThrottled && stakedEntityCount[factory.String()] > 1 {
			manager.logger.Debug("skipping throttled factory", "sender", entry.UserOp.Sender(), "nonce", entry.UserOp.Nonce())
			continue
		}
		if _, ok := senders[entry.UserOp.Sender().String()]; ok {
			manager.logger.Debug("skipping already included sender", "sender", entry.UserOp.Sender(), "nonce", entry.UserOp.Nonce())
			continue
		}

		// re-validate UserOp. no need to check stake, since it cannot be reduced between first and 2nd validation
		validationResult, err := manager.validationManager.ValidateUserOp(entry.UserOp, false)
		if err != nil {
			manager.logger.Debug("failed 2nd validation:", "sender", entry.UserOp.Sender(), "nonce", entry.UserOp.Nonce(), "err", err.Error())
			manager.mempoolManager.RemoveUserOp(entry.UserOp)
			continue
		}
		// TODO: we take UserOp's callGasLimit, even though it will probably require less
		// (but we don't attempt to estimate it to check)
		// which means we could "cram" more UserOps into a bundle.
		userOpGasCost := new(big.Int).Add(validationResult.ReturnInfo.PreOpGas, entry.UserOp.CallGasLimit())
		newTotalGas := new(big.Int).Add(totalGas, userOpGasCost)

		if newTotalGas.Cmp(new(big.Int).SetUint64(manager.maxBundleGas)) > 0 {
			break
		}

		if paymaster != nil {
			_, ok := paymasterDeposit[paymaster.String()]
			if !ok {
				paymasterDeposit[paymaster.String()], err = manager.entryPoint.Caller().BalanceOf(&bind.CallOpts{}, *paymaster)
				if err != nil {
					manager.logger.Error("cannot get balance of paymaster:", "sender", entry.UserOp.Sender(), "nonce", entry.UserOp.Nonce(), "paymaster", paymaster, "err", err.Error())
					continue
				}
			}
			if paymasterDeposit[paymaster.String()].Cmp(validationResult.ReturnInfo.Prefund) < 0 {
				// not enough balance in paymaster to pay for all UserOps
				// (but it passed validation, so it can sponsor them separately
				continue
			}
			stakedEntityCount[paymaster.String()] += 1
			paymasterDeposit[paymaster.String()] = new(big.Int).Sub(paymasterDeposit[paymaster.String()], validationResult.ReturnInfo.Prefund)
		}

		if factory != nil {
			stakedEntityCount[factory.String()] += 1
		}

		senders[entry.UserOp.Sender().String()] = struct{}{}
		bundle = append(bundle, entry.UserOp)
		totalGas = newTotalGas
	}

	return bundle
}

func getAddr(data []byte) *common.Address {
	if len(data) < 20 {
		return nil
	}

	addr := common.BytesToAddress(data[:20])
	return &addr
}
