package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.OpMsg = &MsgEthereumOp{}
)

func (m *MsgEthereumOp) ValidateBasic() error {
	return nil
}

func (m *MsgEthereumOp) Hash() common.Hash {
	return common.Hash{}
}

func (m *MsgEthereumOp) GetEntryPointAddress() *common.Address {
	addr := common.HexToAddress(m.EntryPoint)
	return &addr
}
