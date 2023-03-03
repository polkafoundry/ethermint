package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/evmos/ethermint/eip4337"
	"github.com/evmos/ethermint/x/evm/types"
)

var _ types.OpMsgServer = &Keeper{}

func (k *Keeper) EthereumOp(goCtx context.Context, opMsg *types.MsgEthereumOp) (*types.MsgEthereumOpResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	data, err := k.entryPointContractABI.Pack("simulateValidation", opMsg.Operation)
	if err != nil {
		return &types.MsgEthereumOpResponse{}, err
	}

	args, err := json.Marshal(&types.TransactionArgs{
		To:   opMsg.GetEntryPointAddress(),
		Data: (*hexutil.Bytes)(&data),
	})
	if err != nil {
		return &types.MsgEthereumOpResponse{}, err
	}

	resp, err := k.EthCall(ctx, &types.EthCallRequest{
		Args:    args,
		ChainId: k.ChainID().Int64(),
	})
	if err != nil {
		return &types.MsgEthereumOpResponse{}, err
	}

	validationResult, err := k.unpackValidationResult(resp.Revert())
	if err != nil {
		return &types.MsgEthereumOpResponse{}, err
	}

	_ = validationResult

	return &types.MsgEthereumOpResponse{}, nil
}

func (k *Keeper) unpackValidationResult(data []byte) (*eip4337.ValidationResult, error) {
	e, ok := k.entryPointContractABI.Errors["ValidationResult"]
	if !ok {
		return nil, errors.New("invalid entry point contract")
	}
	if len(data) < 4 {
		return nil, errors.New("invalid data for unpacking")
	}
	if !bytes.Equal(data[:4], e.ID[:4]) {
		return nil, errors.New("invalid data for unpacking")
	}
	unpacked, err := e.Inputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	var ret eip4337.ValidationResult
	err = e.Inputs.Copy(&ret, unpacked)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
