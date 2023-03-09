package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/eip4337"
	"github.com/evmos/ethermint/types"
	"math/big"
)

var (
	_ sdk.OpMsg = &OpMsgEthereum{}
)

func NewOpMsgEthereum(op Operation, chainID *math.Int, entryPoint string) *OpMsgEthereum {
	addr := common.HexToAddress(entryPoint)
	hash := UserOperationHash(op, chainID.BigInt(), &addr)
	return &OpMsgEthereum{
		Operation:  op,
		EntryPoint: entryPoint,
		ChainID:    chainID,
		Hash:       hash.Hex(),
		Bundler:    "",
	}
}

func (m *OpMsgEthereum) ValidateBasic() error {
	if m.Bundler != "" {
		if err := types.ValidateAddress(m.Bundler); err != nil {
			return errorsmod.Wrap(err, "invalid bundler address")
		}
	}

	if err := types.ValidateAddress(m.EntryPoint); err != nil {
		return errorsmod.Wrap(err, "invalid entrypoint address")
	}

	if m.ChainID == nil {
		return errorsmod.Wrapf(errortypes.ErrInvalidRequest, "invalid chain id")
	}

	if err := m.Operation.Validate(); err != nil {
		return err
	}

	hash := UserOperationHash(m.Operation, m.ChainID.BigInt(), m.GetEntryPointAddress()).Hex()
	if m.Hash != hash {
		return errorsmod.Wrapf(errortypes.ErrInvalidRequest, "invalid tx hash %s, expected: %s", m.Hash, hash)
	}

	return nil
}

func (m *OpMsgEthereum) GetEntryPointAddress() *common.Address {
	addr := common.HexToAddress(m.EntryPoint)
	return &addr
}

func (m Operation) Validate() error {
	return nil
}

func UserOperationHash(op Operation, chainID *big.Int, entryPoint *common.Address) common.Hash {
	// FIXME: implement
	return common.Hash{}
}

func (m ReturnInfo) Parse() eip4337.ReturnInfo {
	var preOpGas, preFund *big.Int

	if len(m.PreOpGas) > 0 {
		preOpGas = new(big.Int).SetBytes(m.PreOpGas)
	}

	if len(m.PreFund) > 0 {
		preFund = new(big.Int).SetBytes(m.PreFund)
	}

	return eip4337.ReturnInfo{
		PreOpGas:         preOpGas,
		Prefund:          preFund,
		SigFailed:        m.SigFailed,
		ValidAfter:       m.ValidAfter,
		ValidUntil:       m.ValidUntil,
		PaymasterContext: common.CopyBytes(m.PaymasterContext),
	}
}

func (m StakeInfo) Parse() eip4337.StakeInfo {
	var stake, unstakeDelaySec *big.Int

	if len(m.Stake) > 0 {
		stake = new(big.Int).SetBytes(m.Stake)
	}

	if len(m.UnstakeDelaySec) > 0 {
		unstakeDelaySec = new(big.Int).SetBytes(m.UnstakeDelaySec)
	}

	return eip4337.StakeInfo{
		Stake:           stake,
		UnstakeDelaySec: unstakeDelaySec,
	}
}

func (m AggregatorStakeInfo) Parse() eip4337.AggregatorStakeInfo {
	addr := common.HexToAddress(m.ActualAggregator)
	return eip4337.AggregatorStakeInfo{
		ActualAggregator: &addr,
		StakeInfo:        m.StakeInfo.Parse(),
	}
}

func (m ValidationResult) Parse() eip4337.ValidationResult {
	return eip4337.ValidationResult{
		ReturnInfo:    m.ReturnInfo.Parse(),
		SenderInfo:    m.SenderInfo.Parse(),
		FactoryInfo:   m.FactoryInfo.Parse(),
		PaymasterInfo: m.PaymasterInfo.Parse(),
	}
}

func (m ValidationResultWithAggregation) Parse() eip4337.ValidationResultWithAggregation {
	return eip4337.ValidationResultWithAggregation{
		ReturnInfo:     m.ReturnInfo.Parse(),
		SenderInfo:     m.SenderInfo.Parse(),
		FactoryInfo:    m.FactoryInfo.Parse(),
		PaymasterInfo:  m.PaymasterInfo.Parse(),
		AggregatorInfo: m.AggregatorInfo.Parse(),
	}
}

func FromEIP4337ReturnInfo(returnInfo eip4337.ReturnInfo) ReturnInfo {
	return ReturnInfo{
		PreOpGas:         returnInfo.PreOpGas.Bytes(),
		PreFund:          returnInfo.Prefund.Bytes(),
		SigFailed:        returnInfo.SigFailed,
		ValidAfter:       returnInfo.ValidAfter,
		ValidUntil:       returnInfo.ValidUntil,
		PaymasterContext: common.CopyBytes(returnInfo.PaymasterContext),
	}
}

func FromEIP4337StakeInfo(stakeInfo eip4337.StakeInfo) StakeInfo {
	return StakeInfo{
		Stake:           stakeInfo.Stake.Bytes(),
		UnstakeDelaySec: stakeInfo.UnstakeDelaySec.Bytes(),
	}
}

func FromEIP4337AggregatorStakeInfo(aggregatorStakeInfo eip4337.AggregatorStakeInfo) AggregatorStakeInfo {
	var strAddr string
	if aggregatorStakeInfo.ActualAggregator != nil {
		strAddr = aggregatorStakeInfo.ActualAggregator.Hex()
	}
	return AggregatorStakeInfo{
		ActualAggregator: strAddr,
		StakeInfo:        FromEIP4337StakeInfo(aggregatorStakeInfo.StakeInfo),
	}
}

func FromEIP4337ValidationResult(validationResult eip4337.ValidationResult) ValidationResult {
	return ValidationResult{
		ReturnInfo:    FromEIP4337ReturnInfo(validationResult.ReturnInfo),
		SenderInfo:    FromEIP4337StakeInfo(validationResult.SenderInfo),
		FactoryInfo:   FromEIP4337StakeInfo(validationResult.FactoryInfo),
		PaymasterInfo: FromEIP4337StakeInfo(validationResult.PaymasterInfo),
	}
}

func FromEIP4337ValidationResultWithAggregation(validationResult eip4337.ValidationResultWithAggregation) ValidationResultWithAggregation {
	return ValidationResultWithAggregation{
		ReturnInfo:     FromEIP4337ReturnInfo(validationResult.ReturnInfo),
		SenderInfo:     FromEIP4337StakeInfo(validationResult.SenderInfo),
		FactoryInfo:    FromEIP4337StakeInfo(validationResult.FactoryInfo),
		PaymasterInfo:  FromEIP4337StakeInfo(validationResult.PaymasterInfo),
		AggregatorInfo: FromEIP4337AggregatorStakeInfo(validationResult.AggregatorInfo),
	}
}
