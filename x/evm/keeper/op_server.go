package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/evmos/ethermint/eip4337"
	"github.com/evmos/ethermint/x/evm/types"
)

var _ types.OpServer = &Keeper{}

func (k *Keeper) EthereumOp(goCtx context.Context, opMsg *types.OpMsgEthereum) (*types.OpMsgEthereumResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: check entry point address is valid (in supported list)
	// TODO: other sanity checks. Few of them need to be checked in another places like anteHandlers
	// - either the sender is an existing contract, or the initCode is not empty, but not both
	// - if initCode is not empty, parse its first 20 bytes as a factory address. Record whether the factory is staked
	// - the verificationGasLimit is sufficiently low (<= MAX_VERIFICATION_GAS) and the preVerificationGas is sufficiently high.
	// (enough to pay for the calldata gas cost of serializing the UserOperation plus PRE_VERIFICATION_OVERHEAD_GAS
	// - The paymasterAndData is either empty, or start with the paymaster address, which is a contract that
	// (i) currently has nonempty code on chain,
	// (ii) has a sufficient deposit to pay for the UserOperation, and
	// (iii) is not currently banned.
	// During simulation, the paymaster’s stake is also checked, depending on its storage usage - see reputation, throttling and banning section for details.
	// the callgas is at least the cost of a CALL with non-zero value.
	// the maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client is willing to accept.
	// At the minimum, they are sufficiently high to be included with the current block.basefee
	// - The sender doesn’t have another UserOperation already present in the pool
	// (or it replaces an existing entry with the same sender and nonce, with a higher maxPriorityFeePerGas and an equally increased maxFeePerGas).
	// Only one UserOperation per sender may be included in a single batch.
	// A sender is exempt from this rule and may have multiple UserOperations in the pool and in a batch if it is staked
	// (see reputation, throttling and banning section below), but this exception is of limited use to normal accounts

	data, err := k.entryPointContractABI.Pack("simulateValidation", opMsg.Operation)
	if err != nil {
		return &types.OpMsgEthereumResponse{}, errorsmod.Wrapf(sdkerrors.ErrLogic, "cannot pack data: %v", err)
	}

	args, err := json.Marshal(&types.TransactionArgs{
		To:   opMsg.GetEntryPointAddress(),
		Data: (*hexutil.Bytes)(&data),
	})
	if err != nil {
		return &types.OpMsgEthereumResponse{}, errorsmod.Wrapf(sdkerrors.ErrLogic, "cannot to marshal eth call arguments: %v", err)
	}

	resp, err := k.EthCall(ctx, &types.EthCallRequest{
		Args:    args,
		ChainId: k.ChainID().Int64(),
	})
	if err != nil {
		return &types.OpMsgEthereumResponse{}, errorsmod.Wrapf(sdkerrors.ErrLogic, "cannot execute eth call: %v", err)
	}

	validationResult, err := k.unpackValidationResult(resp.Revert())
	if err != nil {
		return &types.OpMsgEthereumResponse{}, errorsmod.Wrapf(sdkerrors.ErrLogic, "simulate validation did not return a valid result: %v", resp.Ret)
	}

	// classify errors

	return &types.OpMsgEthereumResponse{
		Hash:             opMsg.Hash,
		ValidationResult: types.FromEIP4337ValidationResult(validationResult),
		Ret:              nil,
	}, nil
}

func (k *Keeper) unpackValidationResult(data []byte) (eip4337.ValidationResult, error) {
	var ret eip4337.ValidationResult

	e, ok := k.entryPointContractABI.Errors["ValidationResult"]
	if !ok {
		return ret, errors.New("invalid entry point contract")
	}
	if len(data) < 4 {
		return ret, errors.New("invalid data for unpacking")
	}
	if !bytes.Equal(data[:4], e.ID[:4]) {
		return ret, errors.New("invalid data for unpacking")
	}
	unpacked, err := e.Inputs.Unpack(data)
	if err != nil {
		return ret, err
	}
	err = e.Inputs.Copy(&ret, unpacked)
	return ret, err
}
