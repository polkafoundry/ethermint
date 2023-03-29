package eip4337

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	entrypoint_interface "github.com/evmos/ethermint/eip4337/entrypoint"
	"github.com/evmos/ethermint/eip4337/log"
	"github.com/evmos/ethermint/eip4337/types"
	"github.com/pkg/errors"
)

type ExecutionManager struct {
	logger                   log.Logger
	reputationCronCancel     context.CancelFunc
	autoBundleIntervalCancel context.CancelFunc
	maxMempoolSize           uint64
	autoBundleInterval       uint64
	bundleMtx                *sync.RWMutex
	mtx                      *sync.RWMutex
	entryPoint               entrypoint_interface.IEntryPoint
	reputationManager        IReputationManager
	mempoolManager           IMempoolManager
	bundleManager            *BundleManager
	validationManager        IValidationManager
	eventsManager            IEventsManager
	provider                 IProvider
}

func NewExecutionManager(
	logger log.Logger,
	provider IProvider,
	entryPoint entrypoint_interface.IEntryPoint,
	bundleManager *BundleManager,
	mempoolManager IMempoolManager,
	reputationManager IReputationManager,
	validationManager IValidationManager,
	eventsManager IEventsManager,
) *ExecutionManager {
	return &ExecutionManager{
		logger:            log.EnsureLogger(logger),
		provider:          provider,
		entryPoint:        entryPoint,
		bundleManager:     bundleManager,
		mempoolManager:    mempoolManager,
		reputationManager: reputationManager,
		validationManager: validationManager,
		eventsManager:     eventsManager,

		maxMempoolSize:     0,
		autoBundleInterval: 0,
		bundleMtx:          &sync.RWMutex{},
		mtx:                &sync.RWMutex{},
	}
}

func (manager *ExecutionManager) SendUserOperation(userOpArgs types.UserOperationArgs, entryPoint common.Address) (common.Hash, error) {
	userOp := types.NewUserOperation(userOpArgs)
	err := manager.validationManager.ValidateUserOpBasic(userOp, entryPoint, true, true)
	if err != nil {
		return common.Hash{}, err
	}

	validationResult, err := manager.validationManager.ValidateUserOp(userOp, false)
	if err != nil {
		return common.Hash{}, err
	}

	chainID, err := manager.provider.ChainID()
	if err != nil {
		return common.Hash{}, NewRPCError(ErrorCodeUnknown, err.Error(), nil)
	}

	// if userOp passes the basic validation check, entryPoint is guaranteed to be not nil
	userOpHash := GetUserOpHash(userOp, entryPoint, chainID)
	err = manager.mempoolManager.AddUserOp(
		userOp,
		userOpHash,
		validationResult.ReturnInfo.Prefund,
		validationResult.SenderInfo,
		types.ReferencedCodeHashes{}, // FIXME: not used yet
		validationResult.AggregatorInfo.Address,
	)
	if err != nil {
		return common.Hash{}, err
	}
	go func() { _, _ = manager.attemptBundle(false) }()
	return userOpHash, nil
}

func (manager *ExecutionManager) attemptBundle(force bool) (SendBundleReturn, error) {
	if !force && manager.mempoolManager.Count() < manager.maxMempoolSize {
		return SendBundleReturn{}, ErrNotEnoughOps
	}

	manager.bundleMtx.Lock()
	defer manager.bundleMtx.Unlock()

	ret, err := manager.bundleManager.SendNextBundle()
	if err != nil {
		return SendBundleReturn{}, err
	}
	if manager.maxMempoolSize == 0 {
		err = manager.bundleManager.HandlePastEvents()
		if err != nil {
			manager.logger.Error("bundleManager failed to handle past events", "error", err)
		}
	}
	return ret, nil
}

func (manager *ExecutionManager) EstimateUserOperationGas(userOpArgs types.UserOperationArgs, entryPoint common.Address) (*types.EstimateUserOpGasResult, error) {
	if userOpArgs.PaymasterAndData == nil {
		userOpArgs.PaymasterAndData = &hexutil.Bytes{}
	}
	if userOpArgs.MaxFeePerGas == nil {
		userOpArgs.MaxFeePerGas = (*hexutil.Big)(new(big.Int).SetInt64(0))
	}
	if userOpArgs.MaxPriorityFeePerGas == nil {
		userOpArgs.MaxPriorityFeePerGas = (*hexutil.Big)(new(big.Int).SetInt64(0))
	}
	if userOpArgs.PreVerificationGas == nil {
		userOpArgs.PreVerificationGas = (*hexutil.Big)(new(big.Int).SetInt64(0))
	}
	if userOpArgs.VerificationGasLimit == nil {
		userOpArgs.VerificationGasLimit = (*hexutil.Big)(new(big.Int).SetInt64(1000000))
	}

	userOp := types.NewUserOperation(userOpArgs)
	validationResult, err := manager.validationManager.ValidateUserOp(userOp, false)
	if err != nil {
		return nil, err
	}

	callGasLimit, err := manager.provider.EstimateGas(ethereum.CallMsg{
		From: manager.entryPoint.Address(),
		To:   userOp.Sender(),
		Data: userOp.CallData(),
	})
	if err != nil {
		re, _ := regexp.Compile(`reason="(.*?)"`)
		msg := re.FindString(err.Error())
		if msg == "" {
			msg = "execution reverted"
		}
		return nil, NewRPCError(ErrorCodeUserOperationReverted, msg, nil)
	}

	var validAfter, validUntil *hexutil.Uint64
	if validationResult.ReturnInfo.ValidAfter != 0 {
		validAfter = (*hexutil.Uint64)(&validationResult.ReturnInfo.ValidAfter)
	}
	if validationResult.ReturnInfo.ValidUntil != 0 {
		validUntil = (*hexutil.Uint64)(&validationResult.ReturnInfo.ValidUntil)
	}

	return &types.EstimateUserOpGasResult{
		PreVerificationGas: (*hexutil.Big)(CalcPreVerificationGas(userOp, DefaultGasOverheads())),
		VerificationGas:    (*hexutil.Big)(validationResult.ReturnInfo.PreOpGas),
		ValidAfter:         validAfter,
		ValidUntil:         validUntil,
		CallGasLimit:       (*hexutil.Uint64)(&callGasLimit),
	}, nil
}

func (manager *ExecutionManager) GetUserOperationByHash(userOpHash common.Hash) (*types.UserOperationByHashResponse, error) {
	iterator, err := manager.entryPoint.Filterer().FilterUserOperationEvent(&bind.FilterOpts{}, [][32]byte{userOpHash}, []common.Address{}, []common.Address{})
	if err != nil {
		return nil, NewRPCError(ErrorCodeUnknown, fmt.Sprintf("cannot filter event: %s", err.Error()), nil)
	}
	defer iterator.Close()

	if !iterator.Next() {
		if iterator.Error() != nil {
			return nil, NewRPCError(ErrorCodeUnknown, iterator.Error().Error(), nil)
		}
		return nil, nil
	}

	evt := iterator.Event
	tx, err := manager.provider.GetTransactionByHash(evt.Raw.TxHash)
	if err != nil {
		return nil, NewRPCError(ErrorCodeUnknown, fmt.Sprintf("cannot find tx: %s", err.Error()), nil)
	}

	if tx.To().Hex() != manager.entryPoint.Address().Hex() {
		return nil, NewRPCError(ErrorCodeUnknown, fmt.Sprintf("cannot parse transaction"), nil)
	}

	userOps, _, err := manager.entryPoint.Decoder().DecodeHandleOps(tx.Data())
	if err != nil {
		return nil, NewRPCError(ErrorCodeUnknown, fmt.Sprintf("cannot decode tx data: %s", err.Error()), nil)
	}

	receipt, err := manager.provider.GetTransactionReceipt(evt.Raw.TxHash)
	if err != nil {
		return nil, NewRPCError(ErrorCodeUnknown, fmt.Sprintf("cannot get transaction receipt: %s", err.Error()), nil)
	}

	entryPointAddr := manager.entryPoint.Address()

	for _, userOp := range userOps {
		if userOp.Sender.Hex() == evt.Sender.Hex() && userOp.Nonce.Cmp(evt.Nonce) == 0 {
			return &types.UserOperationByHashResponse{
				UserOperation: types.UserOperationResponse{
					Sender:               &userOp.Sender,
					Nonce:                (*hexutil.Big)(userOp.Nonce),
					InitCode:             (*hexutil.Bytes)(&userOp.InitCode),
					CallData:             (*hexutil.Bytes)(&userOp.CallData),
					CallGasLimit:         (*hexutil.Big)(userOp.CallGasLimit),
					VerificationGasLimit: (*hexutil.Big)(userOp.VerificationGasLimit),
					PreVerificationGas:   (*hexutil.Big)(userOp.PreVerificationGas),
					MaxFeePerGas:         (*hexutil.Big)(userOp.MaxFeePerGas),
					MaxPriorityFeePerGas: (*hexutil.Big)(userOp.MaxPriorityFeePerGas),
					PaymasterAndData:     (*hexutil.Bytes)(&userOp.PaymasterAndData),
					Signature:            (*hexutil.Bytes)(&userOp.Signature),
				},
				EntryPoint:      &entryPointAddr,
				BlockNumber:     (*hexutil.Big)(receipt.BlockNumber),
				BlockHash:       &receipt.BlockHash,
				TransactionHash: &receipt.TxHash,
			}, nil
		}
	}

	return nil, NewRPCError(ErrorCodeUnknown, fmt.Sprintf("cannot find userOp in transaction"), nil)
}

func (manager *ExecutionManager) GetUserOperationReceipt(userOpHash common.Hash) (*types.UserOperationReceipt, error) {
	iterator, err := manager.entryPoint.Filterer().FilterUserOperationEvent(&bind.FilterOpts{}, [][32]byte{userOpHash}, []common.Address{}, []common.Address{})
	if err != nil {
		return nil, NewRPCError(ErrorCodeUnknown, fmt.Sprintf("cannot filter event: %s", err.Error()), nil)
	}
	defer iterator.Close()

	if !iterator.Next() {
		if iterator.Error() != nil {
			return nil, NewRPCError(ErrorCodeUnknown, iterator.Error().Error(), nil)
		}
		return nil, nil
	}

	evt := iterator.Event
	receipt, err := manager.provider.GetTransactionReceipt(evt.Raw.TxHash)
	if err != nil {
		return nil, NewRPCError(ErrorCodeUnknown, fmt.Sprintf("cannot get transaction receipt: %s", err.Error()), nil)
	}

	logs, err := filterLogs(evt, receipt.Logs)
	if err != nil {
		return nil, NewRPCError(ErrorCodeUnknown, fmt.Sprintf("cannot get log from receipt: %s", err.Error()), nil)
	}

	return &types.UserOperationReceipt{
		UserOpHash:    &userOpHash,
		Sender:        &evt.Sender,
		Nonce:         (*hexutil.Big)(evt.Nonce),
		Paymaster:     &evt.Paymaster,
		ActualGasCode: (*hexutil.Big)(evt.ActualGasCost),
		ActualGasUsed: (*hexutil.Big)(evt.ActualGasUsed),
		Success:       evt.Success,
		Reason:        nil,
		Logs:          logs,
		Receipt:       receipt,
	}, nil
}

func (manager *ExecutionManager) SupportedEntryPoints() ([]common.Address, error) {
	return []common.Address{manager.entryPoint.Address()}, nil
}

func filterLogs(userOpEvent *entrypoint_interface.EntryPointUserOperationEvent, logs []*ethtypes.Log) ([]*ethtypes.Log, error) {
	startIndex := -1
	endIndex := -1
	for idx, ethLog := range logs {
		if ethLog == nil {
			continue
		}
		if len(ethLog.Topics) == 0 || len(userOpEvent.Raw.Topics) == 0 {
			continue
		}
		if ethLog.Topics[0].String() == userOpEvent.Raw.Topics[0].String() {
			if len(ethLog.Topics) > 1 && len(userOpEvent.Raw.Topics) > 1 && ethLog.Topics[1].String() == userOpEvent.Raw.Topics[1].String() {
				endIndex = idx
			} else {
				if endIndex == -1 {
					startIndex = idx
				}
			}
		}
	}
	if endIndex == -1 {
		return nil, errors.New("no UserOperationEvent in logs")
	}

	return logs[startIndex+1 : endIndex], nil
}

func (manager *ExecutionManager) SetBundlingInterval(autoBundleInterval uint64, maxMempoolSize uint64) error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	if manager.autoBundleIntervalCancel != nil {
		manager.autoBundleIntervalCancel()
	}

	manager.autoBundleInterval = autoBundleInterval
	manager.maxMempoolSize = maxMempoolSize

	if manager.autoBundleInterval > 0 {
		ctx, cancel := context.WithCancel(context.Background())
		manager.autoBundleIntervalCancel = cancel

		go func(ctx context.Context) {
			timer := time.NewTimer(0)
			for {
				select {
				case <-ctx.Done():
					return
				case <-timer.C:
					start := time.Now()
					_, err := manager.attemptBundle(false)
					if err != nil && err != ErrNotEnoughOps {
						manager.logger.Error("auto bundle failed attempt", "error", err)
					}
					bundleDuration := time.Since(start)
					sleepDuration := time.Duration(manager.autoBundleInterval)*time.Second - bundleDuration
					if sleepDuration < 0 {
						sleepDuration = 0
					}
					timer.Reset(sleepDuration)
				}
			}
		}(ctx)
	}

	return nil
}

/*
- THESE METHODS BELOW IS FOR DEBUG API
*/

func (manager *ExecutionManager) ClearState() error {
	manager.mempoolManager.ClearState()
	manager.reputationManager.ClearState()
	return nil
}

func (manager *ExecutionManager) DumpMempool() ([]types.UserOperationResponse, error) {
	entries := manager.mempoolManager.Dump()
	responses := make([]types.UserOperationResponse, len(entries))
	for idx, entry := range entries {
		responses[idx] = types.ToUserOperationResponse(entry)
	}
	return responses, nil
}

func (manager *ExecutionManager) SendBundleNow() (common.Hash, error) {
	resp, err := manager.attemptBundle(true)
	if err != nil {
		return common.Hash{}, err
	}

	// HandlePastEvents is performed before processing the next bundle.
	// however, in debug mode, we are interested in the side effects
	// (on the mempool) of this "sendBundle" operation
	err = manager.eventsManager.HandlePastEvents()
	if err != nil {
		return common.Hash{}, err
	}

	return resp.TransactionHash, nil
}

func (manager *ExecutionManager) SetBundlingMode(mode string) error {
	if mode == "manual" {
		return manager.SetBundlingInterval(0, 1<<64-1)
	}
	if mode == "auto" {
		return manager.SetBundlingInterval(0, 0)
	}
	return NewRPCError(ErrorCodeInvalidRequest, "mode must be either 'manual' or 'auto'", nil)
}

func (manager *ExecutionManager) SetReputation(reputation types.ReputationArgs) error {
	entries := []ReputationEntry{
		{
			Address:     reputation.Address,
			OpsSeen:     reputation.OpsSeen,
			OpsIncluded: reputation.OpsIncluded,
		},
	}
	manager.reputationManager.SetReputation(entries)

	return nil
}

func (manager *ExecutionManager) DumpReputation() ([]types.ReputationResponse, error) {
	entries := manager.reputationManager.Dump()
	response := make([]types.ReputationResponse, len(entries))
	for idx, entry := range entries {
		response[idx] = types.ReputationResponse{
			Address:     entry.Address,
			OpsSeen:     entry.OpsSeen,
			OpsIncluded: entry.OpsIncluded,
			Status:      string(manager.reputationManager.Status(&entry.Address)),
		}
	}
	return response, nil
}

func (manager *ExecutionManager) Initialize() error {
	// FIXME: do we really need to lock here?
	err := manager.eventsManager.InitEventListener()
	if err != nil {
		return errors.Wrap(err, "failed to init event listener")
	}
	err = manager.eventsManager.InitialHandlePastEvents()
	if err != nil {
		return errors.Wrap(err, "failed to initial handle past events")
	}
	return nil
}
