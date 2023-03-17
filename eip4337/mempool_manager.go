package eip4337

import (
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/eip4337/types"
	"github.com/pkg/errors"
)

type IMempoolManager interface {
	Count() uint64
	GetSortedForInclusion() []MempoolEntry
	AddUserOp(op types.UserOperation, userOpHash common.Hash, prefund *big.Int, senderInfo StakeInfo, referencedContracts types.ReferencedCodeHashes, aggregator *common.Address) error
	RemoveUserOp(op types.UserOperation)
	RemoveUserByHash(userOpHash common.Hash)
	Dump() []types.UserOperation
	ClearState()
}

const MaxMempoolUserOpsPerSender = 4

var (
	_             IMempoolManager = (*MempoolManager)(nil)
	ErrTooManyOps                 = errors.New("sender already has too many operations in mempool")
)

type MempoolManager struct {
	mtx               *sync.RWMutex
	entries           []MempoolEntry
	entryCount        map[string]int
	reputationManager IReputationManager
}

type MempoolEntry struct {
	UserOp              types.UserOperation
	UserOpHash          common.Hash
	Prefund             *big.Int
	ReferencedContracts types.ReferencedCodeHashes
	Aggregator          *common.Address
}

func NewMempoolManager(reputationManager IReputationManager) *MempoolManager {
	return &MempoolManager{
		mtx:               &sync.RWMutex{},
		entries:           make([]MempoolEntry, 0),
		entryCount:        make(map[string]int),
		reputationManager: reputationManager,
	}
}

func (manager *MempoolManager) Count() uint64 {
	return uint64(len(manager.entries))
}

func (manager *MempoolManager) GetSortedForInclusion() []MempoolEntry {
	manager.mtx.RLock()
	defer manager.mtx.RUnlock()
	cpy := append([]MempoolEntry{}, manager.entries...)

	sort.SliceStable(cpy, func(i, j int) bool {
		// TODO: sort in which order ?
		// TODO: need to consult baseFee and maxFeePerGas
		if cpy[j].UserOp.MaxPriorityFeePerGas() == nil {
			return true
		}
		if cpy[i].UserOp.MaxPriorityFeePerGas() == nil {
			return false
		}
		return cpy[i].UserOp.MaxPriorityFeePerGas().Cmp(cpy[j].UserOp.MaxPriorityFeePerGas()) > 0
	})

	return cpy
}

func (manager *MempoolManager) AddUserOp(op types.UserOperation, userOpHash common.Hash, prefund *big.Int, senderInfo StakeInfo, referencedContracts types.ReferencedCodeHashes, aggregator *common.Address) error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()

	entry := MempoolEntry{
		UserOp:              op,
		UserOpHash:          userOpHash,
		Prefund:             prefund,
		ReferencedContracts: referencedContracts,
		Aggregator:          aggregator,
	}

	idx := manager.findBySenderNonce(op.Sender(), op.Nonce())
	if idx != -1 {
		oldEntry := manager.entries[idx]
		oldMaxPriorityFeePerGas := oldEntry.UserOp.MaxPriorityFeePerGas()
		newMaxPriorityFeePerGas := entry.UserOp.MaxPriorityFeePerGas()
		oldMaxFeePerGas := oldEntry.UserOp.MaxFeePerGas()
		newMaxFeePerGas := entry.UserOp.MaxFeePerGas()
		big11 := new(big.Int).SetInt64(11)
		big10 := new(big.Int).SetInt64(10)

		if newMaxPriorityFeePerGas.Cmp(new(big.Int).Div(new(big.Int).Mul(oldMaxPriorityFeePerGas, big11), big10)) < 0 {
			return NewRPCError(
				ErrorCodeInvalidFields,
				fmt.Sprintf(`replacement UserOperation must have higher maxPriorityFeePerGas (old=%s, new=%s)`, oldMaxPriorityFeePerGas.String(), newMaxPriorityFeePerGas.String()),
				nil,
			)
		}

		if newMaxFeePerGas.Cmp(new(big.Int).Div(new(big.Int).Mul(oldMaxFeePerGas, big11), big10)) < 0 {
			return NewRPCError(
				ErrorCodeInvalidFields,
				fmt.Sprintf(`replacement UserOperation must have higher maxFeePerGas (old=%s, new=%s)`, oldMaxFeePerGas.String(), newMaxFeePerGas.String()),
				nil,
			)
		}

		manager.entries[idx] = entry
	} else {
		if manager.entryCount[entry.UserOp.Sender().String()] >= MaxMempoolUserOpsPerSender && manager.reputationManager.CheckStake("account", senderInfo) != nil {
			return ErrTooManyOps
		}
		manager.entryCount[entry.UserOp.Sender().String()] += 1
		manager.entries = append(manager.entries, entry)
	}

	manager.reputationManager.UpdateSeenStatus(aggregator)
	manager.reputationManager.UpdateSeenStatus(getAddr(op.PaymasterAndData()))
	manager.reputationManager.UpdateSeenStatus(getAddr(op.InitCode()))

	return nil
}

// RemoveUserOp remove user operation from mempool by sender and nonce.
// it assumes sender and nonce are not nil
func (manager *MempoolManager) RemoveUserOp(removeOp types.UserOperation) {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	idx := manager.findBySenderNonce(removeOp.Sender(), removeOp.Nonce())
	if idx == -1 {
		return
	}
	op := manager.entries[idx].UserOp
	manager.entries = append(manager.entries[:idx], manager.entries[idx+1:]...)
	count := manager.entryCount[op.Sender().String()] - 1
	if count <= 0 {
		delete(manager.entryCount, op.Sender().String())
	} else {
		manager.entryCount[op.Sender().String()] = count
	}
}

func (manager *MempoolManager) RemoveUserByHash(userOpHash common.Hash) {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	idx := manager.findByHash(userOpHash)
	if idx == -1 {
		return
	}
	op := manager.entries[idx].UserOp
	manager.entries = append(manager.entries[:idx], manager.entries[idx+1:]...)
	count := manager.entryCount[op.Sender().String()] - 1
	if count <= 0 {
		delete(manager.entryCount, op.Sender().String())
	} else {
		manager.entryCount[op.Sender().String()] = count
	}
}

func (manager *MempoolManager) Dump() []types.UserOperation {
	manager.mtx.RLock()
	defer manager.mtx.RUnlock()
	userOps := make([]types.UserOperation, 0)
	for _, entry := range manager.entries {
		userOps = append(userOps, entry.UserOp)
	}
	return userOps
}

func (manager *MempoolManager) ClearState() {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	manager.entries = make([]MempoolEntry, 0)
}

func (manager *MempoolManager) findBySenderNonce(addr *common.Address, nonce *big.Int) int {
	for idx, entry := range manager.entries {
		if entry.UserOp.Sender().String() == addr.String() && entry.UserOp.Nonce().Cmp(nonce) == 0 {
			return idx
		}
	}
	return -1
}

func (manager *MempoolManager) findByHash(hash common.Hash) int {
	for idx, entry := range manager.entries {
		if entry.UserOpHash.String() == hash.String() {
			return idx
		}
	}
	return -1
}
