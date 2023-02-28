package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.Op = &MsgEthereumOp{}
)

func (m *MsgEthereumOp) ValidateBasic() error {
	return nil
}
