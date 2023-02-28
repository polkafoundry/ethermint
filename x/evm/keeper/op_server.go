package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/x/evm/types"
)

var _ types.OpMsgServer = &Keeper{}

func (k *Keeper) EthereumOp(goCtx context.Context, op *types.MsgEthereumOp) (*types.MsgEthereumOpResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	//args, err := json.Marshal(&types.TransactionArgs{
	//	From: &op.EntryPoint,
	//	Data: (*hexutil.Bytes)(&op.Operation.CallData),
	//})
	//
	//_, err = k.EthCall(ctx, &types.EthCallRequest{
	//	Args:            nil,
	//	GasCap:          0,
	//	ProposerAddress: nil,
	//	ChainId:         k.ChainID().Int64(),
	//})
	//if err != nil {
	//	return &types.MsgEthereumOpResponse{}, err
	//}
	return &types.MsgEthereumOpResponse{}, nil
}
