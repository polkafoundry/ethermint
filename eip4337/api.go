package eip4337

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/eip4337/types"
)

type IPublicAPI interface {
	SendUserOperation(userOpArgs types.UserOperationArgs, entryPoint common.Address) (common.Hash, error)
	EstimateUserOperationGas(userOpArgs types.UserOperationArgs, entryPoint common.Address) (*types.EstimateUserOpGasResult, error)
	GetUserOperationReceipt(userOpHash common.Hash) (*types.UserOperationReceipt, error)
	SupportedEntryPoints() ([]common.Address, error)
}

type PublicAPI struct {
	executionManager *ExecutionManager
}

var _ IPublicAPI = (*PublicAPI)(nil)

func NewPublicAPI(executionManager *ExecutionManager) *PublicAPI {
	return &PublicAPI{
		executionManager: executionManager,
	}
}

func (api *PublicAPI) SendUserOperation(userOpArgs types.UserOperationArgs, entryPoint common.Address) (common.Hash, error) {
	return api.executionManager.SendUserOperation(userOpArgs, entryPoint)
}

func (api *PublicAPI) EstimateUserOperationGas(userOpArgs types.UserOperationArgs, entryPoint common.Address) (*types.EstimateUserOpGasResult, error) {
	return api.executionManager.EstimateUserOperationGas(userOpArgs, entryPoint)
}

func (api *PublicAPI) GetUserOperationReceipt(userOpHash common.Hash) (*types.UserOperationReceipt, error) {
	return api.executionManager.GetUserOperationReceipt(userOpHash)
}

func (api *PublicAPI) SupportedEntryPoints() ([]common.Address, error) {
	return api.executionManager.SupportedEntryPoints()
}

// IDebugPublicAPI is public api for bundler debug.
// These methods' names look weird,
// but it is required to make ethereum rpc able to resolve handlers from names
type IDebugPublicAPI interface {
	Bundler_clearState() error
	Bundler_dumpMempool() ([]types.UserOperationResponse, error)
	Bundler_sendBundleNow() (common.Hash, error)
	Bundler_setBundlingMode(mode string) error
	Bundler_setBundlingInterval(interval uint64, maxMempoolSize uint64) error
	Bundler_setReputation(reputation types.ReputationArgs) error
	Bundler_dumpReputation() ([]types.ReputationResponse, error)
}

var _ IDebugPublicAPI = (*DebugPublicAPI)(nil)

type DebugPublicAPI struct {
	executionManager *ExecutionManager
}

func NewDebugAPI(executionManager *ExecutionManager) IDebugPublicAPI {
	return &DebugPublicAPI{
		executionManager: executionManager,
	}
}

func (api *DebugPublicAPI) Bundler_clearState() error {
	return api.executionManager.ClearState()
}

func (api *DebugPublicAPI) Bundler_dumpMempool() ([]types.UserOperationResponse, error) {
	return api.executionManager.DumpMempool()
}

func (api *DebugPublicAPI) Bundler_sendBundleNow() (common.Hash, error) {
	return api.executionManager.SendBundleNow()
}

func (api *DebugPublicAPI) Bundler_setBundlingMode(mode string) error {
	return api.executionManager.SetBundlingMode(mode)
}

func (api *DebugPublicAPI) Bundler_setBundlingInterval(interval uint64, maxMempoolSize uint64) error {
	return api.executionManager.SetBundlingInterval(interval, maxMempoolSize)
}

func (api *DebugPublicAPI) Bundler_setReputation(reputation types.ReputationArgs) error {
	return api.executionManager.SetReputation(reputation)
}

func (api *DebugPublicAPI) Bundler_dumpReputation() ([]types.ReputationResponse, error) {
	return api.executionManager.Bundler_dumpReputation()
}
