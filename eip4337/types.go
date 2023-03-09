package eip4337

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type UserOperation struct {
	evmtypes.Operation
}

func NewInvalidOperationError(err error) *OpError {
	return &OpError{
		Code: -32602,
		Err:  err,
		Data: nil,
	}
}

func NewInternalError(err error) *OpError {
	return &OpError{
		Code: -90909,
		Err:  err,
		Data: nil,
	}
}

func NewErrorWithData(code int, msg string, data interface{}) *OpError {
	return &OpError{
		Code: code,
		Err:  errors.New(msg),
		Data: data,
	}
}

type OpError struct {
	Code int
	Err  error
	Data interface{}
}

type OpErrorData struct {
	Paymaster           string `json:"paymaster,omitempty"`
	ValidUntil          uint64 `json:"validUntil,omitempty"`
	ValidAfter          uint64 `json:"validAfter,omitempty"`
	Aggregator          string `json:"aggregator,omitempty"`
	MinimumStake        uint64 `json:"minimumStake,omitempty"`
	MinimumUnstakeDelay uint64 `json:"minimumUnstakeDelay,omitempty"`
}

func (err *OpError) Error() string {
	return err.Err.Error()
}

func (err *OpError) ErrorCode() int {
	return err.Code
}

func (err *OpError) ErrorData() interface{} {
	return err.Data
}

func (err *OpError) Unwrap() error {
	return err.Err
}

type ValidationResult struct {
	ReturnInfo    ReturnInfo
	SenderInfo    StakeInfo
	FactoryInfo   StakeInfo
	PaymasterInfo StakeInfo
}

type ValidationResultWithAggregation struct {
	ReturnInfo     ReturnInfo
	SenderInfo     StakeInfo
	FactoryInfo    StakeInfo
	PaymasterInfo  StakeInfo
	AggregatorInfo AggregatorStakeInfo
}

type ReturnInfo struct {
	PreOpGas         *big.Int
	Prefund          *big.Int
	SigFailed        bool
	ValidAfter       uint64
	ValidUntil       uint64
	PaymasterContext []byte
}

type StakeInfo struct {
	Stake           *big.Int
	UnstakeDelaySec *big.Int
}

type AggregatorStakeInfo struct {
	ActualAggregator *common.Address
	StakeInfo        StakeInfo
}
