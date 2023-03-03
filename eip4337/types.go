package eip4337

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type UserOperation struct {
}

type OpError struct {
	Code    int
	Message string
	Data    OpErrorData
}

type OpErrorData struct {
	Paymaster           string
	ValidUntil          uint64
	ValidAfter          uint64
	Aggregator          string
	MinimumStake        uint64
	MinimumUnstakeDelay uint64
}

type ValidationResult struct {
	ReturnInfo    ReturnInfo
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
