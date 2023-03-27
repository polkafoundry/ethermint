package eip4337

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/eip4337/log"
	"github.com/evmos/ethermint/rpc/backend"
	rpctypes "github.com/evmos/ethermint/rpc/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/pkg/errors"
)

type IProvider interface {
	SendRawTransaction(txBytes []byte) error
	ChainID() (*big.Int, error)
	GetGasPrice() (*big.Int, error)
	GetBalance(address common.Address) (*big.Int, error)
	GetTransactionCount(address common.Address) (uint64, error)
	GetTransactionReceipt(txHash common.Hash) (*ethtypes.Receipt, error)
	GetTransactionByHash(txHash common.Hash) (*ethtypes.Transaction, error)
	EstimateGas(call ethereum.CallMsg) (uint64, error)
	CurrentHeader() *ethtypes.Header
}

type Provider struct {
	logger  log.Logger
	backend backend.EVMBackend
}

var _ IProvider = (*Provider)(nil)

func NewProvider(logger log.Logger, backend backend.EVMBackend) IProvider {
	return &Provider{
		logger:  log.EnsureLogger(logger),
		backend: backend,
	}
}

func (provider *Provider) SendRawTransaction(txBytes []byte) error {
	_, err := provider.backend.SendRawTransaction(txBytes)
	return err
}

func (provider *Provider) ChainID() (*big.Int, error) {
	id, err := provider.backend.ChainID()
	if err != nil {
		return nil, err
	}
	return (*big.Int)(id), nil
}

func (provider *Provider) GetGasPrice() (*big.Int, error) {
	gasPrice, err := provider.backend.GasPrice()
	if err != nil {
		return nil, err
	}
	return (*big.Int)(gasPrice), err
}

func (provider *Provider) EstimateGas(call ethereum.CallMsg) (uint64, error) {
	bn := rpctypes.EthLatestBlockNumber
	gas, err := provider.backend.EstimateGas(evmtypes.TransactionArgs{
		From:                 &call.From,
		To:                   call.To,
		Gas:                  (*hexutil.Uint64)(&call.Gas),
		GasPrice:             (*hexutil.Big)(call.GasPrice),
		MaxFeePerGas:         (*hexutil.Big)(call.GasFeeCap),
		MaxPriorityFeePerGas: (*hexutil.Big)(call.GasTipCap),
		Value:                (*hexutil.Big)(call.Value),
		Nonce:                nil,
		Data:                 (*hexutil.Bytes)(&call.Data),
		Input:                nil,
		AccessList:           &call.AccessList,
		ChainID:              nil,
	}, &bn)
	if err != nil {
		return 0, err
	}
	return uint64(gas), nil
}

func (provider *Provider) GetBalance(address common.Address) (*big.Int, error) {
	bn := rpctypes.EthLatestBlockNumber
	balance, err := provider.backend.GetBalance(address, rpctypes.BlockNumberOrHash{BlockNumber: &bn})
	if err != nil {
		return nil, err
	}
	return (*big.Int)(balance), err
}

func (provider *Provider) GetTransactionCount(address common.Address) (uint64, error) {
	bn := rpctypes.EthLatestBlockNumber
	count, err := provider.backend.GetTransactionCount(address, bn)
	if err != nil {
		return 0, err
	}
	return uint64(*count), nil
}

func (provider *Provider) GetTransactionReceipt(txHash common.Hash) (*ethtypes.Receipt, error) {
	receiptMap, err := provider.backend.GetTransactionReceipt(txHash)
	if err != nil {
		return nil, err
	}

	if receiptMap == nil {
		return nil, ethereum.NotFound
	}

	receipt := &ethtypes.Receipt{
		Type:              0,
		PostState:         nil,
		Status:            0,
		CumulativeGasUsed: 0,
		Bloom:             ethtypes.Bloom{},
		Logs:              nil,
		TxHash:            txHash,
		ContractAddress:   common.Address{},
		GasUsed:           0,
		BlockHash:         common.Hash{},
		BlockNumber:       nil,
		TransactionIndex:  0,
	}

	if v, ok := receiptMap["type"].(hexutil.Uint); ok {
		receipt.Type = uint8(v)
	}

	if v, ok := receiptMap["status"].(hexutil.Uint); ok {
		receipt.Status = uint64(v)
	}

	if v, ok := receiptMap["cumulativeGasUsed"].(hexutil.Uint64); ok {
		receipt.CumulativeGasUsed = uint64(v)
	}

	if v, ok := receiptMap["logsBloom"].(ethtypes.Bloom); ok {
		receipt.Bloom = v
	}

	if v, ok := receiptMap["logs"].([]*ethtypes.Log); ok {
		receipt.Logs = v
	}

	if v, ok := receiptMap["contractAddress"].(common.Address); ok {
		receipt.ContractAddress = v
	}

	if v, ok := receiptMap["gasUsed"].(hexutil.Uint64); ok {
		receipt.GasUsed = uint64(v)
	}

	if v, ok := receiptMap["blockHash"].(string); ok {
		receipt.BlockHash = common.HexToHash(v)
	}

	if v, ok := receiptMap["blockNumber"].(hexutil.Uint64); ok {
		receipt.BlockNumber = new(big.Int).SetUint64(uint64(v))
	}

	if v, ok := receiptMap["transactionIndex"].(hexutil.Uint64); ok {
		receipt.TransactionIndex = uint(v)
	}

	return receipt, nil
}

func (provider *Provider) GetTransactionByHash(txHash common.Hash) (*ethtypes.Transaction, error) {
	tx, err := provider.backend.GetTransactionByHash(txHash)
	if err != nil {
		return nil, err
	}

	var txData ethtypes.TxData
	switch tx.Type {
	case ethtypes.LegacyTxType:
		txData = &ethtypes.LegacyTx{
			Nonce:    uint64(tx.Nonce),
			GasPrice: (*big.Int)(tx.GasPrice),
			Gas:      uint64(tx.Gas),
			To:       tx.To,
			Value:    (*big.Int)(tx.Value),
			Data:     tx.Input,
			V:        (*big.Int)(tx.V),
			R:        (*big.Int)(tx.R),
			S:        (*big.Int)(tx.S),
		}
	case ethtypes.AccessListTxType:
		txData = &ethtypes.AccessListTx{
			ChainID:    (*big.Int)(tx.ChainID),
			Nonce:      uint64(tx.Nonce),
			GasPrice:   (*big.Int)(tx.GasPrice),
			Gas:        uint64(tx.Gas),
			To:         tx.To,
			Value:      (*big.Int)(tx.Value),
			Data:       tx.Input,
			AccessList: *tx.Accesses,
			V:          (*big.Int)(tx.V),
			R:          (*big.Int)(tx.R),
			S:          (*big.Int)(tx.S),
		}
	case ethtypes.DynamicFeeTxType:
		txData = &ethtypes.DynamicFeeTx{
			ChainID:    (*big.Int)(tx.ChainID),
			Nonce:      uint64(tx.Nonce),
			GasTipCap:  (*big.Int)(tx.GasTipCap),
			GasFeeCap:  (*big.Int)(tx.GasFeeCap),
			Gas:        uint64(tx.Gas),
			To:         tx.To,
			Value:      (*big.Int)(tx.Value),
			Data:       tx.Input,
			AccessList: *tx.Accesses,
			V:          (*big.Int)(tx.V),
			R:          (*big.Int)(tx.R),
			S:          (*big.Int)(tx.S),
		}
	default:
		return nil, errors.New("invalid transaction type")
	}

	return ethtypes.NewTx(txData), nil
}

func (provider *Provider) CurrentHeader() *ethtypes.Header {
	return provider.backend.CurrentHeader()
}
