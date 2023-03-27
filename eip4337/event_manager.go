package eip4337

import (
	"bytes"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	entrypoint_interface "github.com/evmos/ethermint/eip4337/entrypoint"
	"github.com/pkg/errors"
)

type IEventsManager interface {
	InitEventListener() error
	HandlePastEvents() error
	InitialHandlePastEvents() error
	Stop()
}

type EventsManager struct {
	lastBlock         uint64
	provider          IProvider
	entryPoint        entrypoint_interface.IEntryPoint
	mempoolManager    IMempoolManager
	reputationManager IReputationManager

	userOperationEventSubscription         event.Subscription
	accountDeployedSubscription            event.Subscription
	signatureAggregatorChangedSubscription event.Subscription

	eventAggregator       *common.Address
	eventAggregatorTxHash *common.Hash

	mtx *sync.RWMutex
}

var _ IEventsManager = (*EventsManager)(nil)

func NewEventsManager(provider IProvider, entryPoint entrypoint_interface.IEntryPoint, mempoolManager IMempoolManager, reputationManager IReputationManager) *EventsManager {
	return &EventsManager{
		lastBlock:         0, // FIXME: since the rpc limit number of logs returned, correct this
		provider:          provider,
		entryPoint:        entryPoint,
		mempoolManager:    mempoolManager,
		reputationManager: reputationManager,
		mtx:               &sync.RWMutex{},
	}
}

func (manager *EventsManager) InitEventListener() error {
	err := manager.initUserOperationEventEventListener()
	if err != nil {
		return errors.Wrap(err, "cannot initiate UserOperationEvent listener")
	}
	//err = manager.initAccountDeployedEventListener()
	//if err != nil {
	//	return errors.Wrap(err, "cannot initiate AccountDeployed listener")
	//}
	//err = manager.initSignatureAggregatorChanged()
	//if err != nil {
	//	return errors.Wrap(err, "cannot initiate SignatureAggregatorChanged listener")
	//}
	return nil
}

func (manager *EventsManager) initUserOperationEventEventListener() error {
	ch := make(chan *entrypoint_interface.EntryPointUserOperationEvent)
	subscription, err := manager.entryPoint.Filterer().WatchUserOperationEvent(&bind.WatchOpts{}, ch, [][32]byte{}, []common.Address{}, []common.Address{})
	if err != nil {
		return err
	}
	manager.userOperationEventSubscription = subscription
	go func() {
		for evt := range ch {
			func() {
				manager.mtx.Lock()
				defer manager.mtx.Unlock()
				manager.handleEvent(evt)
			}()
		}
	}()
	return nil
}

func (manager *EventsManager) initAccountDeployedEventListener() error {
	ch := make(chan *entrypoint_interface.EntryPointAccountDeployed)
	subscription, err := manager.entryPoint.Filterer().WatchAccountDeployed(&bind.WatchOpts{}, ch, [][32]byte{}, []common.Address{})
	if err != nil {
		return err
	}
	manager.userOperationEventSubscription = subscription
	go func() {
		for evt := range ch {
			func() {
				manager.mtx.Lock()
				defer manager.mtx.Unlock()
				manager.handleEvent(evt)
			}()
		}
	}()
	return nil
}

func (manager *EventsManager) initSignatureAggregatorChanged() error {
	ch := make(chan *entrypoint_interface.EntryPointSignatureAggregatorChanged)
	subscription, err := manager.entryPoint.Filterer().WatchSignatureAggregatorChanged(&bind.WatchOpts{}, ch, []common.Address{})
	if err != nil {
		return err
	}
	manager.userOperationEventSubscription = subscription
	go func() {
		for evt := range ch {
			func() {
				manager.mtx.Lock()
				defer manager.mtx.Unlock()
				manager.handleEvent(evt)
			}()
		}
	}()
	return nil
}

var zeroAddress = common.Address{}

// handleEvent requires caller to have lock
func (manager *EventsManager) handleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *entrypoint_interface.EntryPointUserOperationEvent:
		manager.handleUserOperationEvent(v)
		manager.lastBlock = v.Raw.BlockNumber + 1
	case *entrypoint_interface.EntryPointAccountDeployed:
		manager.handleAccountDeployedEvent(v)
		manager.lastBlock = v.Raw.BlockNumber + 1
	case *entrypoint_interface.EntryPointSignatureAggregatorChanged:
		manager.handleSignatureAggregatorChangedEvent(v)
		manager.lastBlock = v.Raw.BlockNumber + 1
	}
}

func (manager *EventsManager) handleUserOperationEvent(evt *entrypoint_interface.EntryPointUserOperationEvent) {
	userOpHash := common.Hash(evt.UserOpHash)
	manager.mempoolManager.RemoveUserByHash(userOpHash)
	manager.includedAddress(evt.Sender)
	manager.includedAddress(evt.Paymaster)

	aggregatorAddr := manager.getEventAggregator(evt)
	if aggregatorAddr != nil {
		manager.includedAddress(*aggregatorAddr)
	}
}

func (manager *EventsManager) handleAccountDeployedEvent(evt *entrypoint_interface.EntryPointAccountDeployed) {
	manager.includedAddress(evt.Paymaster)
}

func (manager *EventsManager) handleSignatureAggregatorChangedEvent(evt *entrypoint_interface.EntryPointSignatureAggregatorChanged) {
	manager.eventAggregator = &evt.Aggregator
	manager.eventAggregatorTxHash = &evt.Raw.TxHash
}

func (manager *EventsManager) includedAddress(addr common.Address) {
	if !bytes.Equal(addr[:], zeroAddress[:]) {
		manager.reputationManager.UpdateIncludedStatus(&addr)
	}
}

// getEventAggregator requires caller to lock
func (manager *EventsManager) getEventAggregator(evt *entrypoint_interface.EntryPointUserOperationEvent) *common.Address {
	if manager.eventAggregatorTxHash == nil || !bytes.Equal(evt.Raw.TxHash[:], manager.eventAggregatorTxHash[:]) {
		manager.eventAggregator = nil
		manager.eventAggregatorTxHash = &evt.Raw.TxHash
	}
	return manager.eventAggregator
}

func (manager *EventsManager) HandlePastEvents() error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	currentHeader := manager.provider.CurrentHeader()
	if currentHeader == nil {
		return errors.New("currentHeader is nil")
	}

	iterator, err := manager.entryPoint.FilterLogs(&bind.FilterOpts{Start: manager.lastBlock})
	if err != nil {
		return err
	}

	for iterator.Next() {
		manager.handleEvent(iterator.Event)
	}

	manager.lastBlock = currentHeader.Number.Uint64()

	return nil
}

func (manager *EventsManager) InitialHandlePastEvents() error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	currentHeader := manager.provider.CurrentHeader()
	if currentHeader == nil {
		return errors.New("currentHeader is nil")
	}

	manager.lastBlock = currentHeader.Number.Uint64() - 86400
	if manager.lastBlock < 0 {
		manager.lastBlock = 0
	}

	iterator, err := manager.entryPoint.FilterLogs(&bind.FilterOpts{Start: manager.lastBlock})
	if err != nil {
		return err
	}

	for iterator.Next() {
		manager.handleEvent(iterator.Event)
	}

	return nil
}

func (manager *EventsManager) Stop() {
	if manager.userOperationEventSubscription != nil {
		manager.userOperationEventSubscription.Unsubscribe()
	}
	if manager.accountDeployedSubscription != nil {
		manager.accountDeployedSubscription.Unsubscribe()
	}
	if manager.signatureAggregatorChangedSubscription != nil {
		manager.signatureAggregatorChangedSubscription.Unsubscribe()
	}
}
