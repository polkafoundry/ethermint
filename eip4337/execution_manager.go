package eip4337

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/evmos/ethermint/eip4337/log"
	"github.com/evmos/ethermint/eip4337/types"
	"github.com/pkg/errors"
)

type ExecutionManager struct {
	logger                   log.Logger
	reputationCronCancel     context.CancelFunc
	autoBundleIntervalCancel context.CancelFunc
	maxMempoolSize           uint64
	autoBundleInterval       time.Duration
	chainID                  *big.Int
	bundleMtx                *sync.RWMutex
	mtx                      *sync.RWMutex
	reputationManager        IReputationManager
	mempoolManager           IMempoolManager
	bundleManager            *BundleManager
	validationManager        IValidationManager
}

func NewExecutionManager(
	logger log.Logger,
	reputationManager IReputationManager,
	mempoolManager IMempoolManager,
	bundleManager *BundleManager,
	validationManager IValidationManager,
) *ExecutionManager {
	return &ExecutionManager{
		logger:             log.EnsureLogger(logger),
		maxMempoolSize:     0,
		autoBundleInterval: 0,
		bundleMtx:          &sync.RWMutex{},
		mtx:                &sync.RWMutex{},
		reputationManager:  reputationManager,
		mempoolManager:     mempoolManager,
		bundleManager:      bundleManager,
		validationManager:  validationManager,
	}
}

func (manager *ExecutionManager) SendUserOperation(userOpArgs types.UserOperationArgs, entryPoint *common.Address) (common.Hash, error) {
	userOp := types.NewUserOperation(userOpArgs)
	err := manager.validationManager.ValidateUserOpBasic(userOp, entryPoint, true, true)
	if err != nil {
		return common.Hash{}, err
	}

	validationResult, err := manager.validationManager.ValidateUserOp(userOp, false)
	if err != nil {
		return common.Hash{}, err
	}

	// if userOp passes the basic validation check, entryPoint is guaranteed to be not nil
	userOpHash := GetUserOpHash(userOp, *entryPoint, manager.chainID)
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
		return SendBundleReturn{}, errors.New("not enough ops to bundle")
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

func (manager *ExecutionManager) SetAutoBundle(autoBundleInterval hexutil.Uint64, maxMempoolSize hexutil.Uint64) error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	if manager.autoBundleIntervalCancel != nil {
		manager.autoBundleIntervalCancel()
	}

	manager.autoBundleInterval = time.Duration(autoBundleInterval) * time.Second
	manager.maxMempoolSize = uint64(maxMempoolSize)

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
					if err != nil {
						manager.logger.Error("auto bundle failed attempt", "error", err)
					}
					bundleDuration := time.Since(start)
					sleepDuration := manager.autoBundleInterval - bundleDuration
					if sleepDuration > 0 {
						sleepDuration = 0
					}
					timer.Reset(sleepDuration)
				}
			}
		}(ctx)
	}

	return nil
}
