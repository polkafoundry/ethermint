package eip4337

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/eip4337/log"
)

type ReputationStatus string

const (
	ReputationStatusOk        ReputationStatus = "ok"
	ReputationStatusThrottled ReputationStatus = "throttled"
	ReputationStatusBanned    ReputationStatus = "banned"
)

type IReputationManager interface {
	AddWhitelist(addresses ...common.Address)
	AddBlacklist(addresses ...common.Address)
	HourlyCron()
	Status(addr *common.Address) ReputationStatus
	CheckStake(title string, stakeInfo StakeInfo) error
	CrashedHandleOps(addr *common.Address)
	UpdateSeenStatus(addr *common.Address)
	UpdateIncludedStatus(addr *common.Address)
	ClearState()
	Dump() []ReputationEntry
	SetReputation(reputations []ReputationEntry) []ReputationEntry
}

var _ IReputationManager = (*ReputationManager)(nil)

type ReputationManager struct {
	mtx             *sync.RWMutex
	logger          log.Logger
	params          ReputationParams
	minStake        *big.Int
	minUnstakeDelay uint64
	entries         map[string]ReputationEntry
	blackList       map[string]struct{}
	whiteList       map[string]struct{}
}

func NewReputationManager(logger log.Logger, params ReputationParams, minStake *big.Int, minUnstakeDelay uint64) *ReputationManager {
	return &ReputationManager{
		mtx:             &sync.RWMutex{},
		logger:          log.EnsureLogger(logger),
		params:          params,
		minStake:        minStake,
		minUnstakeDelay: minUnstakeDelay,
		entries:         make(map[string]ReputationEntry),
		blackList:       make(map[string]struct{}),
		whiteList:       make(map[string]struct{}),
	}
}

type ReputationParams struct {
	MinInclusionDenominator int64
	ThrottlingSlack         int64
	BanSlack                int64
}

func DefaultBundlerReputationParams() ReputationParams {
	return ReputationParams{
		MinInclusionDenominator: 10,
		ThrottlingSlack:         10,
		BanSlack:                50,
	}
}

func DefaultNonBundlerReputationParams() ReputationParams {
	return ReputationParams{
		MinInclusionDenominator: 100,
		ThrottlingSlack:         10,
		BanSlack:                50,
	}
}

type ReputationEntry struct {
	Address     common.Address
	OpsSeen     int64
	OpsIncluded int64
}

func (manager *ReputationManager) AddWhitelist(addresses ...common.Address) {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	for _, addr := range addresses {
		manager.whiteList[addr.String()] = struct{}{}
	}
}

func (manager *ReputationManager) AddBlacklist(addresses ...common.Address) {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	for _, addr := range addresses {
		manager.blackList[addr.String()] = struct{}{}
	}
}

func (manager *ReputationManager) HourlyCron() {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	for key, entry := range manager.entries {
		entry.OpsSeen = entry.OpsSeen * 23 / 24
		entry.OpsIncluded = entry.OpsIncluded * 23 / 24
		if entry.OpsSeen == 0 && entry.OpsIncluded == 0 {
			delete(manager.entries, key)
		}
		manager.entries[key] = entry
	}
}

func (manager *ReputationManager) Status(addr *common.Address) ReputationStatus {
	manager.mtx.RLock()
	defer manager.mtx.RUnlock()

	if addr == nil || manager.isWhitelisted(*addr) {
		return ReputationStatusOk
	}

	if manager.isBlacklisted(*addr) {
		return ReputationStatusBanned
	}

	entry, ok := manager.entries[addr.String()]
	if !ok {
		return ReputationStatusOk
	}

	minExpectedIncluded := entry.OpsIncluded / manager.params.MinInclusionDenominator
	switch {
	case minExpectedIncluded <= entry.OpsIncluded+manager.params.ThrottlingSlack:
		return ReputationStatusOk
	case minExpectedIncluded <= entry.OpsIncluded+manager.params.BanSlack:
		return ReputationStatusThrottled
	default:
		return ReputationStatusBanned
	}
}

func (manager *ReputationManager) CheckStake(title string, stakeInfo StakeInfo) error {
	manager.mtx.RLock()
	defer manager.mtx.RUnlock()

	if stakeInfo.Address == nil || manager.isWhitelisted(*stakeInfo.Address) {
		return nil
	}

	if manager.Status(stakeInfo.Address) == ReputationStatusBanned {
		return NewRPCError(
			ErrorCodeReputation,
			fmt.Sprintf("%s %s is banned", title, stakeInfo.Address.String()),
			map[string]string{title: stakeInfo.Address.String()},
		)
	}

	if stakeInfo.Stake.Cmp(manager.minStake) < 0 {
		return NewRPCError(
			ErrorCodeInsufficientStake,
			fmt.Sprintf("%s %s stake %s is too low (min=%s)", title, stakeInfo.Address.String(), stakeInfo.Stake.String(), manager.minStake.String()),
			nil,
		)
	}

	if stakeInfo.UnstakeDelaySec.Cmp(new(big.Int).SetInt64(int64(manager.minUnstakeDelay))) < 0 {
		return NewRPCError(
			ErrorCodeInsufficientStake,
			fmt.Sprintf("%s %s unstake delay %s is too low (min=%d)", title, stakeInfo.Address.String(), stakeInfo.UnstakeDelaySec.String(), manager.minUnstakeDelay),
			nil,
		)
	}

	return nil
}

func (manager *ReputationManager) CrashedHandleOps(addr *common.Address) {
	if addr == nil {
		return
	}

	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	entry, ok := manager.entries[addr.String()]
	if !ok {
		entry = ReputationEntry{
			Address:     *addr,
			OpsSeen:     0,
			OpsIncluded: 0,
		}
	}

	entry.OpsSeen = 100
	entry.OpsIncluded = 0
	manager.entries[addr.String()] = entry
}

func (manager *ReputationManager) UpdateSeenStatus(addr *common.Address) {
	if addr == nil {
		return
	}

	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	entry, ok := manager.entries[addr.String()]
	if !ok {
		entry = ReputationEntry{
			Address:     *addr,
			OpsSeen:     0,
			OpsIncluded: 0,
		}
	}

	entry.OpsSeen += 1
	manager.entries[addr.String()] = entry
}

func (manager *ReputationManager) UpdateIncludedStatus(addr *common.Address) {
	if addr == nil {
		return
	}

	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	entry, ok := manager.entries[addr.String()]
	if !ok {
		entry = ReputationEntry{
			Address:     *addr,
			OpsSeen:     0,
			OpsIncluded: 0,
		}
	}

	entry.OpsIncluded += 1
	manager.entries[addr.String()] = entry
}

func (manager *ReputationManager) ClearState() {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	manager.entries = make(map[string]ReputationEntry)
}

func (manager *ReputationManager) Dump() []ReputationEntry {
	manager.mtx.RUnlock()
	defer manager.mtx.RUnlock()
	entries := make([]ReputationEntry, 0)
	for _, entry := range manager.entries {
		entries = append(entries, entry)
	}
	return entries
}

func (manager *ReputationManager) SetReputation(reputations []ReputationEntry) []ReputationEntry {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	for _, entry := range reputations {
		manager.entries[entry.Address.String()] = entry
	}

	entries := make([]ReputationEntry, 0)
	for _, entry := range manager.entries {
		entries = append(entries, entry)
	}

	return entries
}

func (manager *ReputationManager) isWhitelisted(address common.Address) bool {
	_, found := manager.whiteList[address.String()]
	return found
}

func (manager *ReputationManager) isBlacklisted(address common.Address) bool {
	_, found := manager.blackList[address.String()]
	return found
}
