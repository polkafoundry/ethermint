package eip4337

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	entrypoint_interface "github.com/evmos/ethermint/eip4337/entrypoint"
	"github.com/evmos/ethermint/eip4337/types"
)

type IValidationManager interface {
	ValidateUserOp(op types.UserOperation, checkStakes bool) (ValidationResult, error)
	ValidateUserOpBasic(op types.UserOperation, entryPointAddress common.Address, requireSignature bool, requireGasParams bool) error
}

var _ IValidationManager = (*ValidationManager)(nil)

type ValidationManager struct {
	entryPoint        entrypoint_interface.IEntryPoint
	reputationManager IReputationManager
	unsafe            bool
}

func NewValidationManager(
	entryPoint entrypoint_interface.IEntryPoint,
	reputationManager IReputationManager,
	unsafe bool,
) IValidationManager {
	return &ValidationManager{
		entryPoint:        entryPoint,
		reputationManager: reputationManager,
		unsafe:            unsafe,
	}
}

type ValidationResult struct {
	ReturnInfo     ReturnInfo
	SenderInfo     StakeInfo
	FactoryInfo    StakeInfo
	PaymasterInfo  StakeInfo
	AggregatorInfo StakeInfo
}

type ReturnInfo struct {
	PreOpGas         *big.Int
	Prefund          *big.Int
	SigFailed        bool
	ValidAfter       uint64
	ValidUntil       uint64
	PaymasterContext []byte
}

func NewReturnInfo(returnInfo entrypoint_interface.ReturnInfo) ReturnInfo {
	ret := ReturnInfo{
		PreOpGas:         nil,
		Prefund:          nil,
		SigFailed:        returnInfo.SigFailed,
		ValidAfter:       returnInfo.ValidAfter,
		ValidUntil:       returnInfo.ValidUntil,
		PaymasterContext: common.CopyBytes(returnInfo.PaymasterContext),
	}
	if returnInfo.PreOpGas != nil {
		ret.PreOpGas = new(big.Int).Set(returnInfo.PreOpGas)
	}
	if returnInfo.Prefund != nil {
		ret.Prefund = new(big.Int).Set(returnInfo.Prefund)
	}
	return ret
}

type StakeInfo struct {
	Address         *common.Address
	Stake           *big.Int
	UnstakeDelaySec *big.Int
}

func NewStakeInfo(addr *common.Address, stakeInfo entrypoint_interface.StakeInfo) StakeInfo {
	ret := StakeInfo{
		Address:         addr,
		Stake:           nil,
		UnstakeDelaySec: nil,
	}
	if stakeInfo.Stake != nil {
		ret.Stake = new(big.Int).Set(stakeInfo.Stake)
	}
	if stakeInfo.UnstakeDelaySec != nil {
		ret.UnstakeDelaySec = new(big.Int).Set(stakeInfo.UnstakeDelaySec)
	}
	return ret
}

type AggregatorStakeInfo struct {
	ActualAggregator *common.Address
	StakeInfo        StakeInfo
}

func (manager *ValidationManager) ValidateUserOp(op types.UserOperation, checkStakes bool) (ValidationResult, error) {
	if manager.unsafe {
		// FIXME: implement me
		return ValidationResult{}, NewRPCError(ErrorCodeUnknown, "unsafe validation manager not implemented", nil)
	}

	res, err := manager.callSimulateValidation(op)
	if err != nil {
		return ValidationResult{}, err
	}

	if res.ReturnInfo.SigFailed {
		return ValidationResult{}, NewRPCError(
			ErrorCodeInvalidSignature,
			"invalid UserOp signature or paymaster signature",
			nil,
		)
	}

	// FIXME: use config instead?
	now := uint64(time.Now().Unix())
	if res.ReturnInfo.ValidAfter > now || res.ReturnInfo.ValidUntil < now+30 {
		// FIXME: The data field SHOULD contain a paymaster value, if this error was triggered by the paymaster
		return ValidationResult{}, NewRPCError(
			ErrorCodeExpiresShortly,
			"userOp is not valid yet or expires too soon",
			map[string]interface{}{
				"validAfter": res.ReturnInfo.ValidAfter,
				"validUntil": res.ReturnInfo.ValidUntil,
			},
		)
	}

	if res.AggregatorInfo.Address != nil {
		err = manager.reputationManager.CheckStake("aggregator", res.AggregatorInfo)
		if err != nil {
			return ValidationResult{}, err
		}

		return ValidationResult{}, NewRPCError(
			ErrorCodeUnsupportedSignatureAggregator,
			"currently not supporting aggregator",
			nil,
		)
	}

	return res, nil
}

func (manager *ValidationManager) callSimulateValidation(op types.UserOperation) (ValidationResult, error) {
	callErr := manager.entryPoint.Caller().SimulateValidation(&bind.CallOpts{}, ToABIUserOperation(op))
	if callErr == nil {
		return ValidationResult{}, NewRPCError(ErrorCodeUnknown, "invalid response, simulateValidation call must revert", nil)
	}

	validationResult, err := manager.entryPoint.ErrorDecoder().DecodeValidationResult(callErr)
	if err == nil {
		return ValidationResult{
			ReturnInfo:     NewReturnInfo(validationResult.ReturnInfo),
			SenderInfo:     NewStakeInfo(op.Sender(), validationResult.SenderInfo),
			FactoryInfo:    NewStakeInfo(getAddr(op.InitCode()), validationResult.FactoryInfo),
			PaymasterInfo:  NewStakeInfo(getAddr(op.PaymasterAndData()), validationResult.PaymasterInfo),
			AggregatorInfo: StakeInfo{},
		}, nil
	}

	validationResultWithAggregation, err := manager.entryPoint.ErrorDecoder().DecodeValidationResultWithAggregation(callErr)
	if err == nil {
		return ValidationResult{
			ReturnInfo:     NewReturnInfo(validationResult.ReturnInfo),
			SenderInfo:     NewStakeInfo(op.Sender(), validationResult.SenderInfo),
			FactoryInfo:    NewStakeInfo(getAddr(op.InitCode()), validationResult.FactoryInfo),
			PaymasterInfo:  NewStakeInfo(getAddr(op.PaymasterAndData()), validationResult.PaymasterInfo),
			AggregatorInfo: NewStakeInfo(validationResultWithAggregation.AggregatorInfo.ActualAggregator, validationResultWithAggregation.AggregatorInfo.StakeInfo),
		}, nil
	}

	failedOp, err := manager.entryPoint.ErrorDecoder().DecodeFailedOp(callErr)
	if err == nil {
		if strings.HasPrefix(failedOp.Reason, "AA3") {
			return ValidationResult{}, NewRPCError(
				ErrorCodeSimulatePaymasterValidation,
				fmt.Sprintf("paymaster validation failed: %s", failedOp.Reason),
				map[string]string{"paymaster": getAddr(op.PaymasterAndData()).String()},
			)
		}
		return ValidationResult{}, NewRPCError(
			ErrorCodeSimulateValidation,
			fmt.Sprintf("account validation failed: %s", failedOp.Reason),
			nil,
		)
	}

	msg, err := manager.entryPoint.ErrorDecoder().DecodeString(callErr)
	if err == nil {
		return ValidationResult{}, NewRPCError(
			ErrorCodeSimulateValidation,
			fmt.Sprintf("account validation failed: %s", msg),
			nil,
		)
	}

	return ValidationResult{}, NewRPCError(
		ErrorCodeUnknown,
		err.Error(),
		nil,
	)
}

func (manager *ValidationManager) ValidateUserOpBasic(op types.UserOperation, entryPointAddress common.Address, requireSignature bool, requireGasParams bool) error {
	if entryPointAddress.String() != manager.entryPoint.Address().String() {
		return NewRPCError(
			ErrorCodeInvalidFields,
			fmt.Sprintf("the entry point at %s is not supported. this bundler uses %s", entryPointAddress.String(), manager.entryPoint.Address().String()),
			nil,
		)
	}

	if op.Sender() == nil {
		return NewRPCError(ErrorCodeInvalidFields, "missing sender field", nil)
	}

	if op.Nonce() == nil {
		return NewRPCError(ErrorCodeInvalidFields, "missing nonce field", nil)
	}

	if op.InitCode() == nil {
		return NewRPCError(ErrorCodeInvalidFields, "missing initCode field", nil)
	}

	if op.CallData() == nil {
		return NewRPCError(ErrorCodeInvalidFields, "missing callData field", nil)
	}

	if op.PaymasterAndData() == nil {
		return NewRPCError(ErrorCodeInvalidFields, "missing paymasterAndData field", nil)
	}

	if requireSignature && op.Signature() == nil {
		return NewRPCError(ErrorCodeInvalidFields, "missing signature field", nil)
	}

	if requireGasParams {
		if op.PreVerificationGas() == nil {
			return NewRPCError(ErrorCodeInvalidFields, "missing preVerificationGas field", nil)
		}

		if op.VerificationGasLimit() == nil {
			return NewRPCError(ErrorCodeInvalidFields, "missing verificationGasLimit field", nil)
		}

		if op.CallGasLimit() == nil {
			return NewRPCError(ErrorCodeInvalidFields, "missing callGasLimit field", nil)
		}

		if op.MaxFeePerGas() == nil {
			return NewRPCError(ErrorCodeInvalidFields, "missing maxFeePerGas field", nil)
		}

		if op.MaxPriorityFeePerGas() == nil {
			return NewRPCError(ErrorCodeInvalidFields, "missing maxPriorityFeePerGas field", nil)
		}
	}

	if len(op.PaymasterAndData()) != 0 && len(op.PaymasterAndData()) < 20 {
		return NewRPCError(ErrorCodeInvalidFields, "paymasterAndData: must contain at least an address", nil)
	}

	if len(op.InitCode()) != 0 && len(op.InitCode()) < 20 {
		return NewRPCError(ErrorCodeInvalidFields, "initCode: must contain at least an address", nil)
	}

	preVerificationGas := CalcPreVerificationGas(op, DefaultGasOverheads())
	if op.PreVerificationGas().Cmp(preVerificationGas) < 0 {
		return NewRPCError(
			ErrorCodeInvalidFields,
			fmt.Sprintf("preVerificationGas too low: expected at least %s", preVerificationGas.String()),
			nil,
		)
	}

	return nil
}

type GasOverheads struct {
	Fixed         int64
	PerUserOp     int64
	PerUserOpWord int64
	ZeroByte      int64
	NonZeroByte   int64
	BundleSize    int
	SigSize       int
}

func DefaultGasOverheads() GasOverheads {
	return GasOverheads{
		Fixed:         21000,
		PerUserOp:     18300,
		PerUserOpWord: 4,
		ZeroByte:      4,
		NonZeroByte:   16,
		BundleSize:    1,
		SigSize:       65,
	}
}

func CalcPreVerificationGas(op types.UserOperation, overheads GasOverheads) *big.Int {
	if op.PreVerificationGas() == nil {
		op = op.WithPreVerificationGas(big.NewInt(21000))
	}
	if op.Signature() == nil {
		op = op.WithSignature(newDummyBytesSliceWithValue(overheads.SigSize, 1))
	}

	packed := PackUserOp(op, false)
	lengthInWord := int64((len(packed) + 31) / 32)
	callDataCost := int64(0)
	for _, x := range packed {
		if x == 0 {
			callDataCost += overheads.ZeroByte
		} else {
			callDataCost += overheads.NonZeroByte
		}
	}

	ret := callDataCost + int64(math.Round(float64(overheads.Fixed)/float64(overheads.BundleSize))) + overheads.PerUserOp + overheads.PerUserOpWord*lengthInWord
	return big.NewInt(ret)
}
