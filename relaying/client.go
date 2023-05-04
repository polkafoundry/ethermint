package relaying

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethfilters "github.com/ethereum/go-ethereum/eth/filters"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/evmos/ethermint/rpc/backend"
	"github.com/evmos/ethermint/rpc/ethereum/pubsub"
	"github.com/evmos/ethermint/rpc/namespaces/ethereum/eth/filters"
	rpctypes "github.com/evmos/ethermint/rpc/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func NewLocalClient(logger log.Logger, eventSystem *filters.EventSystem, backend *backend.Backend) *LocalClient {
	return &LocalClient{
		events:  eventSystem,
		logger:  logger,
		backend: backend,
	}
}

var _ bind.ContractBackend = (*LocalClient)(nil)

type LocalClient struct {
	events  *filters.EventSystem
	logger  log.Logger
	backend *backend.Backend
}

func (client *LocalClient) CodeAt(_ context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	bn := toBlockNumberArg(blockNumber)
	return client.backend.GetCode(contract, rpctypes.BlockNumberOrHash{BlockNumber: &bn})
}

func (client *LocalClient) CallContract(_ context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	bn := toBlockNumberArg(blockNumber)
	resp, err := client.backend.DoCall(evmtypes.TransactionArgs{
		From:                 addrToAddrPtr(call.From),
		To:                   call.To,
		Gas:                  toHexUtilUint64(call.Gas),
		GasPrice:             (*hexutil.Big)(call.GasPrice),
		MaxFeePerGas:         (*hexutil.Big)(call.GasFeeCap),
		MaxPriorityFeePerGas: (*hexutil.Big)(call.GasTipCap),
		Value:                (*hexutil.Big)(call.Value),
		Nonce:                nil,
		Data:                 (*hexutil.Bytes)(&call.Data),
		Input:                nil,
		AccessList:           &call.AccessList,
		ChainID:              nil,
	}, bn)
	if err != nil {
		return []byte{}, err
	}
	return resp.Ret, nil
}

func (client *LocalClient) HeaderByNumber(_ context.Context, blockNumber *big.Int) (*ethtypes.Header, error) {
	bn := toBlockNumberArg(blockNumber)
	return client.backend.HeaderByNumber(bn)
}

func (client *LocalClient) PendingCodeAt(_ context.Context, account common.Address) ([]byte, error) {
	bn := rpctypes.EthPendingBlockNumber
	return client.backend.GetCode(account, rpctypes.BlockNumberOrHash{BlockNumber: &bn})
}

func (client *LocalClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	bn := rpctypes.EthPendingBlockNumber
	nonce, err := client.backend.GetTransactionCount(account, bn)
	if err != nil {
		return 0, err
	}
	return uint64(*nonce), nil
}

func (client *LocalClient) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	gasPrice, err := client.backend.GasPrice()
	if err != nil {
		return nil, err
	}
	return (*big.Int)(gasPrice), nil
}

func (client *LocalClient) SuggestGasTipCap(_ context.Context) (*big.Int, error) {
	head := client.backend.CurrentHeader()
	return client.backend.SuggestGasTipCap(head.BaseFee)
}

func (client *LocalClient) EstimateGas(_ context.Context, call ethereum.CallMsg) (uint64, error) {
	bn := rpctypes.EthPendingBlockNumber
	gas, err := client.backend.EstimateGas(evmtypes.TransactionArgs{
		From:                 addrToAddrPtr(call.From),
		To:                   call.To,
		Gas:                  toHexUtilUint64(call.Gas),
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

func (client *LocalClient) SendTransaction(_ context.Context, tx *ethtypes.Transaction) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = client.backend.SendRawTransaction(data)
	return err
}

func (client *LocalClient) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]ethtypes.Log, error) {
	// copy from rpc/namespaces/ethereum/eth/filters/api.go
	var filter *filters.Filter
	crit := ethfilters.FilterCriteria(query)
	if crit.BlockHash != nil {
		// Block filter requested, construct a single-shot filter
		filter = filters.NewBlockFilter(client.logger, client.backend, crit)
	} else {
		// Convert the RPC block numbers into internal representations
		begin := ethrpc.LatestBlockNumber.Int64()
		if crit.FromBlock != nil {
			begin = crit.FromBlock.Int64()
		}
		end := ethrpc.LatestBlockNumber.Int64()
		if crit.ToBlock != nil {
			end = crit.ToBlock.Int64()
		}
		// Construct the range filter
		filter = filters.NewRangeFilter(client.logger, client.backend, begin, end, crit.Addresses, crit.Topics)
	}

	// Run the filter and return all the logs
	logs, err := filter.Logs(ctx, math.MaxInt, math.MaxInt64)
	if err != nil {
		return nil, err
	}

	returnLogs := make([]ethtypes.Log, 0)
	for _, l := range logs {
		returnLogs = append(returnLogs, *l)
	}

	return returnLogs, nil
}

type subscription struct {
	unsubFn pubsub.UnsubscribeFunc
	err     chan error
}

func (s subscription) Unsubscribe() {
	s.unsubFn()
}

func (s subscription) Err() <-chan error {
	return s.err
}

var _ ethereum.Subscription = (*subscription)(nil)

func (client *LocalClient) SubscribeFilterLogs(_ context.Context, query ethereum.FilterQuery, ch chan<- ethtypes.Log) (ethereum.Subscription, error) {
	crit := ethfilters.FilterCriteria(query)

	sub, unsubFn, err := client.events.SubscribeLogs(crit)
	if err != nil {
		client.logger.Error("failed to subscribe logs", "error", err.Error())
		return nil, err
	}

	subID := ethrpc.NewID()

	s := subscription{
		unsubFn: unsubFn,
		err:     make(chan error, 1),
	}

	go func() {
		eventCh := sub.Event()
		errCh := sub.Err()
		for {
			select {
			case event, ok := <-eventCh:
				if !ok {
					return
				}

				dataTx, ok := event.Data.(tmtypes.EventDataTx)
				if !ok {
					client.logger.Debug("event data type mismatch", "type", fmt.Sprintf("%T", event.Data))
					continue
				}

				txResponse, err := evmtypes.DecodeTxResponse(dataTx.TxResult.Result.Data)
				if err != nil {
					client.logger.Error("failed to decode tx response", "error", err.Error())
					s.err <- err
					return
				}

				logs := filters.FilterLogs(evmtypes.LogsToEthereum(txResponse.Logs), crit.FromBlock, crit.ToBlock, crit.Addresses, crit.Topics)
				if len(logs) == 0 {
					continue
				}

				for _, ethLog := range logs {
					ch <- *ethLog
				}
			case err, ok := <-errCh:
				if !ok {
					s.err <- err
					return
				}
				client.logger.Debug("dropping Logs WebSocket subscription", "subscription-id", subID, "error", err.Error())
			}
		}
	}()

	return s, nil
}
