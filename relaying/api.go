package relaying

import "github.com/ethereum/go-ethereum/common"

type IPublicAPI interface {
	SendRelayTransaction(args RelayTransactionArgs) (common.Hash, error)
	GetRelayReceipt(txHash common.Hash) (RelayReceipt, error)
	GetRefundTokens() []common.Address
	GetRefundAddresses() []common.Address
}

type PublicAPI struct {
	relayer IRelayer
}

var _ IPublicAPI = (*PublicAPI)(nil)

func NewPublicAPI(relayer IRelayer) IPublicAPI {
	return &PublicAPI{relayer: relayer}
}

func (api *PublicAPI) SendRelayTransaction(args RelayTransactionArgs) (common.Hash, error) {
	return api.relayer.SendRelayTransaction(args)
}

func (api *PublicAPI) GetRelayReceipt(txHash common.Hash) (RelayReceipt, error) {
	return api.relayer.GetRelayReceipt(txHash)
}

func (api *PublicAPI) GetRefundTokens() []common.Address {
	return api.relayer.GetRefundTokens()
}

func (api *PublicAPI) GetRefundAddresses() []common.Address {
	return api.relayer.GetRefundAddresses()
}
