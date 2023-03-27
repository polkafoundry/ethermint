package entrypoint_interface

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

type IEntryPointUserOpsPerAggregator struct {
	UserOps    []UserOperation
	Aggregator common.Address
	Signature  []byte
}

type IStakeManagerDepositInfo struct {
	Deposit         *big.Int
	Staked          bool
	Stake           *big.Int
	UnstakeDelaySec uint32
	WithdrawTime    *big.Int
}

type UserOperation struct {
	Sender               common.Address
	Nonce                *big.Int
	InitCode             []byte
	CallData             []byte
	CallGasLimit         *big.Int
	VerificationGasLimit *big.Int
	PreVerificationGas   *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	PaymasterAndData     []byte
	Signature            []byte
}

type FailedOp struct {
	OpIndex *big.Int
	Reason  string
}

type ValidationResult struct {
	ReturnInfo    ReturnInfo
	SenderInfo    StakeInfo
	FactoryInfo   StakeInfo
	PaymasterInfo StakeInfo
}

type ValidationResultWithAggregator struct {
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
	ValidAfter       *big.Int
	ValidUntil       *big.Int
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

var EntryPointMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"preOpGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"},{\"internalType\":\"uint48\",\"name\":\"validAfter\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"validUntil\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"targetSuccess\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"targetResult\",\"type\":\"bytes\"}],\"name\":\"ExecutionResult\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"opIndex\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"FailedOp\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderAddressResult\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"SignatureValidationFailed\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"preOpGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prefund\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"sigFailed\",\"type\":\"bool\"},{\"internalType\":\"uint48\",\"name\":\"validAfter\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"validUntil\",\"type\":\"uint48\"},{\"internalType\":\"bytes\",\"name\":\"paymasterContext\",\"type\":\"bytes\"}],\"internalType\":\"structIEntryPoint.ReturnInfo\",\"name\":\"returnInfo\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeDelaySec\",\"type\":\"uint256\"}],\"internalType\":\"structIStakeManager.StakeInfo\",\"name\":\"senderInfo\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeDelaySec\",\"type\":\"uint256\"}],\"internalType\":\"structIStakeManager.StakeInfo\",\"name\":\"factoryInfo\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeDelaySec\",\"type\":\"uint256\"}],\"internalType\":\"structIStakeManager.StakeInfo\",\"name\":\"paymasterInfo\",\"type\":\"tuple\"}],\"name\":\"ValidationResult\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"preOpGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prefund\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"sigFailed\",\"type\":\"bool\"},{\"internalType\":\"uint48\",\"name\":\"validAfter\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"validUntil\",\"type\":\"uint48\"},{\"internalType\":\"bytes\",\"name\":\"paymasterContext\",\"type\":\"bytes\"}],\"internalType\":\"structIEntryPoint.ReturnInfo\",\"name\":\"returnInfo\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeDelaySec\",\"type\":\"uint256\"}],\"internalType\":\"structIStakeManager.StakeInfo\",\"name\":\"senderInfo\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeDelaySec\",\"type\":\"uint256\"}],\"internalType\":\"structIStakeManager.StakeInfo\",\"name\":\"factoryInfo\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeDelaySec\",\"type\":\"uint256\"}],\"internalType\":\"structIStakeManager.StakeInfo\",\"name\":\"paymasterInfo\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeDelaySec\",\"type\":\"uint256\"}],\"internalType\":\"structIStakeManager.StakeInfo\",\"name\":\"stakeInfo\",\"type\":\"tuple\"}],\"internalType\":\"structIEntryPoint.AggregatorStakeInfo\",\"name\":\"aggregatorInfo\",\"type\":\"tuple\"}],\"name\":\"ValidationResultWithAggregation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"paymaster\",\"type\":\"address\"}],\"name\":\"AccountDeployed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalDeposit\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"SignatureAggregatorChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalStaked\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unstakeDelaySec\",\"type\":\"uint256\"}],\"name\":\"StakeLocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdrawTime\",\"type\":\"uint256\"}],\"name\":\"StakeUnlocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"withdrawAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"StakeWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"paymaster\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"actualGasCost\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"actualGasUsed\",\"type\":\"uint256\"}],\"name\":\"UserOperationEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"revertReason\",\"type\":\"bytes\"}],\"name\":\"UserOperationRevertReason\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"withdrawAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_unstakeDelaySec\",\"type\":\"uint32\"}],\"name\":\"addStake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"depositTo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getDepositInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint112\",\"name\":\"deposit\",\"type\":\"uint112\"},{\"internalType\":\"bool\",\"name\":\"staked\",\"type\":\"bool\"},{\"internalType\":\"uint112\",\"name\":\"stake\",\"type\":\"uint112\"},{\"internalType\":\"uint32\",\"name\":\"unstakeDelaySec\",\"type\":\"uint32\"},{\"internalType\":\"uint48\",\"name\":\"withdrawTime\",\"type\":\"uint48\"}],\"internalType\":\"structIStakeManager.DepositInfo\",\"name\":\"info\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"name\":\"getSenderAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"}],\"name\":\"getUserOpHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation[]\",\"name\":\"userOps\",\"type\":\"tuple[]\"},{\"internalType\":\"contractIAggregator\",\"name\":\"aggregator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structIEntryPoint.UserOpsPerAggregator[]\",\"name\":\"opsPerAggregator\",\"type\":\"tuple[]\"},{\"internalType\":\"addresspayable\",\"name\":\"beneficiary\",\"type\":\"address\"}],\"name\":\"handleAggregatedOps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation[]\",\"name\":\"ops\",\"type\":\"tuple[]\"},{\"internalType\":\"addresspayable\",\"name\":\"beneficiary\",\"type\":\"address\"}],\"name\":\"handleOps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"op\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"targetCallData\",\"type\":\"bytes\"}],\"name\":\"simulateHandleOp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"}],\"name\":\"simulateValidation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlockStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"withdrawAddress\",\"type\":\"address\"}],\"name\":\"withdrawStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"withdrawAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"withdrawAmount\",\"type\":\"uint256\"}],\"name\":\"withdrawTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var EntryPointABI = EntryPointMetaData.ABI

var (
	_ IEntryPointCaller     = (*EntryPointCaller)(nil)
	_ IEntryPointTransactor = (*EntryPointTransactor)(nil)
	_ IEntryPointFilterer   = (*EntryPointFilterer)(nil)
	_ IEntryPointDecoder    = (*EntryPointDecoder)(nil)
)

type IEntryPointCaller interface {
	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)
	GetDepositInfo(opts *bind.CallOpts, account common.Address) (IStakeManagerDepositInfo, error)
	GetUserOpHash(opts *bind.CallOpts, userOp UserOperation) ([32]byte, error)
	SimulateValidation(opts *bind.CallOpts, userOp UserOperation) error
	HandleOps(opts *bind.CallOpts, userOps []UserOperation, beneficiary common.Address) error
}

type IEntryPointTransactor interface {
	AddStake(opts *bind.TransactOpts, _unstakeDelaySec uint32) (*types.Transaction, error)
	DepositTo(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)
	GetSenderAddress(opts *bind.TransactOpts, initCode []byte) (*types.Transaction, error)
	HandleAggregatedOps(opts *bind.TransactOpts, opsPerAggregator []IEntryPointUserOpsPerAggregator, beneficiary common.Address) (*types.Transaction, error)
	HandleOps(opts *bind.TransactOpts, ops []UserOperation, beneficiary common.Address) (*types.Transaction, error)
	SimulateHandleOp(opts *bind.TransactOpts, op UserOperation, target common.Address, targetCallData []byte) (*types.Transaction, error)
	SimulateValidation(opts *bind.TransactOpts, userOp UserOperation) (*types.Transaction, error)
	UnlockStake(opts *bind.TransactOpts) (*types.Transaction, error)
	WithdrawStake(opts *bind.TransactOpts, withdrawAddress common.Address) (*types.Transaction, error)
	WithdrawTo(opts *bind.TransactOpts, withdrawAddress common.Address, withdrawAmount *big.Int) (*types.Transaction, error)
}

type IEntryPointFilterer interface {
	FilterAccountDeployed(opts *bind.FilterOpts, userOpHash [][32]byte, sender []common.Address) (*EntryPointAccountDeployedIterator, error)
	WatchAccountDeployed(opts *bind.WatchOpts, sink chan<- *EntryPointAccountDeployed, userOpHash [][32]byte, sender []common.Address) (event.Subscription, error)
	ParseAccountDeployed(log types.Log) (*EntryPointAccountDeployed, error)
	FilterDeposited(opts *bind.FilterOpts, account []common.Address) (*EntryPointDepositedIterator, error)
	WatchDeposited(opts *bind.WatchOpts, sink chan<- *EntryPointDeposited, account []common.Address) (event.Subscription, error)
	ParseDeposited(log types.Log) (*EntryPointDeposited, error)
	FilterSignatureAggregatorChanged(opts *bind.FilterOpts, aggregator []common.Address) (*EntryPointSignatureAggregatorChangedIterator, error)
	WatchSignatureAggregatorChanged(opts *bind.WatchOpts, sink chan<- *EntryPointSignatureAggregatorChanged, aggregator []common.Address) (event.Subscription, error)
	ParseSignatureAggregatorChanged(log types.Log) (*EntryPointSignatureAggregatorChanged, error)
	FilterStakeLocked(opts *bind.FilterOpts, account []common.Address) (*EntryPointStakeLockedIterator, error)
	WatchStakeLocked(opts *bind.WatchOpts, sink chan<- *EntryPointStakeLocked, account []common.Address) (event.Subscription, error)
	ParseStakeLocked(log types.Log) (*EntryPointStakeLocked, error)
	FilterStakeUnlocked(opts *bind.FilterOpts, account []common.Address) (*EntryPointStakeUnlockedIterator, error)
	WatchStakeUnlocked(opts *bind.WatchOpts, sink chan<- *EntryPointStakeUnlocked, account []common.Address) (event.Subscription, error)
	ParseStakeUnlocked(log types.Log) (*EntryPointStakeUnlocked, error)
	FilterStakeWithdrawn(opts *bind.FilterOpts, account []common.Address) (*EntryPointStakeWithdrawnIterator, error)
	WatchStakeWithdrawn(opts *bind.WatchOpts, sink chan<- *EntryPointStakeWithdrawn, account []common.Address) (event.Subscription, error)
	ParseStakeWithdrawn(log types.Log) (*EntryPointStakeWithdrawn, error)
	FilterUserOperationEvent(opts *bind.FilterOpts, userOpHash [][32]byte, sender []common.Address, paymaster []common.Address) (*EntryPointUserOperationEventIterator, error)
	WatchUserOperationEvent(opts *bind.WatchOpts, sink chan<- *EntryPointUserOperationEvent, userOpHash [][32]byte, sender []common.Address, paymaster []common.Address) (event.Subscription, error)
	ParseUserOperationEvent(log types.Log) (*EntryPointUserOperationEvent, error)
	FilterUserOperationRevertReason(opts *bind.FilterOpts, userOpHash [][32]byte, sender []common.Address) (*EntryPointUserOperationRevertReasonIterator, error)
	WatchUserOperationRevertReason(opts *bind.WatchOpts, sink chan<- *EntryPointUserOperationRevertReason, userOpHash [][32]byte, sender []common.Address) (event.Subscription, error)
	ParseUserOperationRevertReason(log types.Log) (*EntryPointUserOperationRevertReason, error)
	FilterWithdrawn(opts *bind.FilterOpts, account []common.Address) (*EntryPointWithdrawnIterator, error)
	WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *EntryPointWithdrawn, account []common.Address) (event.Subscription, error)
	ParseWithdrawn(log types.Log) (*EntryPointWithdrawn, error)
}

type IEntryPointDecoder interface {
	DecodeErrorFailedOp(error) (FailedOp, error)
	DecodeErrorValidationResult(error) (ValidationResult, error)
	DecodeErrorValidationResultWithAggregation(error) (ValidationResultWithAggregator, error)
	DecodeErrorString(error) (string, error)
	DecodeHandleOps(data []byte) ([]UserOperation, common.Address, error)
}

type IEntryPoint interface {
	Address() common.Address
	Caller() IEntryPointCaller
	Transactor() IEntryPointTransactor
	Filterer() IEntryPointFilterer
	Decoder() IEntryPointDecoder
	FilterLogs(opts *bind.FilterOpts, query ...[]interface{}) (*EntryPointEventIterator, error)
}

type EntryPointCaller struct {
	contract *bind.BoundContract
}

func (_EntryPoint *EntryPointCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EntryPoint.contract.Call(opts, &out, "balanceOf", account)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

func (_EntryPoint *EntryPointCaller) GetDepositInfo(opts *bind.CallOpts, account common.Address) (IStakeManagerDepositInfo, error) {
	var out []interface{}
	err := _EntryPoint.contract.Call(opts, &out, "getDepositInfo", account)
	if err != nil {
		return *new(IStakeManagerDepositInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IStakeManagerDepositInfo)).(*IStakeManagerDepositInfo)

	return out0, err
}

func (_EntryPoint *EntryPointCaller) GetUserOpHash(opts *bind.CallOpts, userOp UserOperation) ([32]byte, error) {
	var out []interface{}
	err := _EntryPoint.contract.Call(opts, &out, "getUserOpHash", userOp)
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

func (_EntryPoint *EntryPointCaller) SimulateValidation(opts *bind.CallOpts, userOp UserOperation) error {
	var out []interface{}
	err := _EntryPoint.contract.Call(opts, &out, "simulateValidation", userOp)
	return err
}

func (_EntryPoint *EntryPointCaller) HandleOps(opts *bind.CallOpts, userOps []UserOperation, beneficiary common.Address) error {
	var out []interface{}
	err := _EntryPoint.contract.Call(opts, &out, "handleOps", userOps, beneficiary)
	return err
}

type EntryPointTransactor struct {
	contract *bind.BoundContract
}

func (_EntryPoint *EntryPointTransactor) AddStake(opts *bind.TransactOpts, _unstakeDelaySec uint32) (*types.Transaction, error) {
	return _EntryPoint.contract.Transact(opts, "addStake", _unstakeDelaySec)
}

func (_EntryPoint *EntryPointTransactor) DepositTo(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _EntryPoint.contract.Transact(opts, "depositTo", account)
}

func (_EntryPoint *EntryPointTransactor) GetSenderAddress(opts *bind.TransactOpts, initCode []byte) (*types.Transaction, error) {
	return _EntryPoint.contract.Transact(opts, "getSenderAddress", initCode)
}

func (_EntryPoint *EntryPointTransactor) HandleAggregatedOps(opts *bind.TransactOpts, opsPerAggregator []IEntryPointUserOpsPerAggregator, beneficiary common.Address) (*types.Transaction, error) {
	return _EntryPoint.contract.Transact(opts, "handleAggregatedOps", opsPerAggregator, beneficiary)
}

func (_EntryPoint *EntryPointTransactor) HandleOps(opts *bind.TransactOpts, ops []UserOperation, beneficiary common.Address) (*types.Transaction, error) {
	return _EntryPoint.contract.Transact(opts, "handleOps", ops, beneficiary)
}

func (_EntryPoint *EntryPointTransactor) SimulateHandleOp(opts *bind.TransactOpts, op UserOperation, target common.Address, targetCallData []byte) (*types.Transaction, error) {
	return _EntryPoint.contract.Transact(opts, "simulateHandleOp", op, target, targetCallData)
}

func (_EntryPoint *EntryPointTransactor) SimulateValidation(opts *bind.TransactOpts, userOp UserOperation) (*types.Transaction, error) {
	return _EntryPoint.contract.Transact(opts, "simulateValidation", userOp)
}

func (_EntryPoint *EntryPointTransactor) UnlockStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EntryPoint.contract.Transact(opts, "unlockStake")
}

func (_EntryPoint *EntryPointTransactor) WithdrawStake(opts *bind.TransactOpts, withdrawAddress common.Address) (*types.Transaction, error) {
	return _EntryPoint.contract.Transact(opts, "withdrawStake", withdrawAddress)
}

func (_EntryPoint *EntryPointTransactor) WithdrawTo(opts *bind.TransactOpts, withdrawAddress common.Address, withdrawAmount *big.Int) (*types.Transaction, error) {
	return _EntryPoint.contract.Transact(opts, "withdrawTo", withdrawAddress, withdrawAmount)
}

type EntryPointFilterer struct {
	contract *bind.BoundContract
}

func (_EntryPoint *EntryPointFilterer) FilterAccountDeployed(opts *bind.FilterOpts, userOpHash [][32]byte, sender []common.Address) (*EntryPointAccountDeployedIterator, error) {
	var userOpHashRule []interface{}
	for _, userOpHashItem := range userOpHash {
		userOpHashRule = append(userOpHashRule, userOpHashItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EntryPoint.contract.FilterLogs(opts, "AccountDeployed", userOpHashRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EntryPointAccountDeployedIterator{contract: _EntryPoint.contract, event: "AccountDeployed", logs: logs, sub: sub}, nil
}

func (_EntryPoint *EntryPointFilterer) WatchAccountDeployed(opts *bind.WatchOpts, sink chan<- *EntryPointAccountDeployed, userOpHash [][32]byte, sender []common.Address) (event.Subscription, error) {
	var userOpHashRule []interface{}
	for _, userOpHashItem := range userOpHash {
		userOpHashRule = append(userOpHashRule, userOpHashItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EntryPoint.contract.WatchLogs(opts, "AccountDeployed", userOpHashRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EntryPointAccountDeployed)
				if err := _EntryPoint.contract.UnpackLog(event, "AccountDeployed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EntryPoint *EntryPointFilterer) ParseAccountDeployed(log types.Log) (*EntryPointAccountDeployed, error) {
	event := new(EntryPointAccountDeployed)
	if err := _EntryPoint.contract.UnpackLog(event, "AccountDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_EntryPoint *EntryPointFilterer) FilterDeposited(opts *bind.FilterOpts, account []common.Address) (*EntryPointDepositedIterator, error) {
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EntryPoint.contract.FilterLogs(opts, "Deposited", accountRule)
	if err != nil {
		return nil, err
	}
	return &EntryPointDepositedIterator{contract: _EntryPoint.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

func (_EntryPoint *EntryPointFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *EntryPointDeposited, account []common.Address) (event.Subscription, error) {
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EntryPoint.contract.WatchLogs(opts, "Deposited", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EntryPointDeposited)
				if err := _EntryPoint.contract.UnpackLog(event, "Deposited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EntryPoint *EntryPointFilterer) ParseDeposited(log types.Log) (*EntryPointDeposited, error) {
	event := new(EntryPointDeposited)
	if err := _EntryPoint.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_EntryPoint *EntryPointFilterer) FilterSignatureAggregatorChanged(opts *bind.FilterOpts, aggregator []common.Address) (*EntryPointSignatureAggregatorChangedIterator, error) {
	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _EntryPoint.contract.FilterLogs(opts, "SignatureAggregatorChanged", aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &EntryPointSignatureAggregatorChangedIterator{contract: _EntryPoint.contract, event: "SignatureAggregatorChanged", logs: logs, sub: sub}, nil
}

func (_EntryPoint *EntryPointFilterer) WatchSignatureAggregatorChanged(opts *bind.WatchOpts, sink chan<- *EntryPointSignatureAggregatorChanged, aggregator []common.Address) (event.Subscription, error) {
	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _EntryPoint.contract.WatchLogs(opts, "SignatureAggregatorChanged", aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EntryPointSignatureAggregatorChanged)
				if err := _EntryPoint.contract.UnpackLog(event, "SignatureAggregatorChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EntryPoint *EntryPointFilterer) ParseSignatureAggregatorChanged(log types.Log) (*EntryPointSignatureAggregatorChanged, error) {
	event := new(EntryPointSignatureAggregatorChanged)
	if err := _EntryPoint.contract.UnpackLog(event, "SignatureAggregatorChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_EntryPoint *EntryPointFilterer) FilterStakeLocked(opts *bind.FilterOpts, account []common.Address) (*EntryPointStakeLockedIterator, error) {
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EntryPoint.contract.FilterLogs(opts, "StakeLocked", accountRule)
	if err != nil {
		return nil, err
	}
	return &EntryPointStakeLockedIterator{contract: _EntryPoint.contract, event: "StakeLocked", logs: logs, sub: sub}, nil
}

func (_EntryPoint *EntryPointFilterer) WatchStakeLocked(opts *bind.WatchOpts, sink chan<- *EntryPointStakeLocked, account []common.Address) (event.Subscription, error) {
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EntryPoint.contract.WatchLogs(opts, "StakeLocked", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EntryPointStakeLocked)
				if err := _EntryPoint.contract.UnpackLog(event, "StakeLocked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EntryPoint *EntryPointFilterer) ParseStakeLocked(log types.Log) (*EntryPointStakeLocked, error) {
	event := new(EntryPointStakeLocked)
	if err := _EntryPoint.contract.UnpackLog(event, "StakeLocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_EntryPoint *EntryPointFilterer) FilterStakeUnlocked(opts *bind.FilterOpts, account []common.Address) (*EntryPointStakeUnlockedIterator, error) {
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EntryPoint.contract.FilterLogs(opts, "StakeUnlocked", accountRule)
	if err != nil {
		return nil, err
	}
	return &EntryPointStakeUnlockedIterator{contract: _EntryPoint.contract, event: "StakeUnlocked", logs: logs, sub: sub}, nil
}

func (_EntryPoint *EntryPointFilterer) WatchStakeUnlocked(opts *bind.WatchOpts, sink chan<- *EntryPointStakeUnlocked, account []common.Address) (event.Subscription, error) {
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EntryPoint.contract.WatchLogs(opts, "StakeUnlocked", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EntryPointStakeUnlocked)
				if err := _EntryPoint.contract.UnpackLog(event, "StakeUnlocked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EntryPoint *EntryPointFilterer) ParseStakeUnlocked(log types.Log) (*EntryPointStakeUnlocked, error) {
	event := new(EntryPointStakeUnlocked)
	if err := _EntryPoint.contract.UnpackLog(event, "StakeUnlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_EntryPoint *EntryPointFilterer) FilterStakeWithdrawn(opts *bind.FilterOpts, account []common.Address) (*EntryPointStakeWithdrawnIterator, error) {
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EntryPoint.contract.FilterLogs(opts, "StakeWithdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return &EntryPointStakeWithdrawnIterator{contract: _EntryPoint.contract, event: "StakeWithdrawn", logs: logs, sub: sub}, nil
}

func (_EntryPoint *EntryPointFilterer) WatchStakeWithdrawn(opts *bind.WatchOpts, sink chan<- *EntryPointStakeWithdrawn, account []common.Address) (event.Subscription, error) {
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EntryPoint.contract.WatchLogs(opts, "StakeWithdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EntryPointStakeWithdrawn)
				if err := _EntryPoint.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EntryPoint *EntryPointFilterer) ParseStakeWithdrawn(log types.Log) (*EntryPointStakeWithdrawn, error) {
	event := new(EntryPointStakeWithdrawn)
	if err := _EntryPoint.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_EntryPoint *EntryPointFilterer) FilterUserOperationEvent(opts *bind.FilterOpts, userOpHash [][32]byte, sender []common.Address, paymaster []common.Address) (*EntryPointUserOperationEventIterator, error) {
	var userOpHashRule []interface{}
	for _, userOpHashItem := range userOpHash {
		userOpHashRule = append(userOpHashRule, userOpHashItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var paymasterRule []interface{}
	for _, paymasterItem := range paymaster {
		paymasterRule = append(paymasterRule, paymasterItem)
	}

	logs, sub, err := _EntryPoint.contract.FilterLogs(opts, "UserOperationEvent", userOpHashRule, senderRule, paymasterRule)
	if err != nil {
		return nil, err
	}
	return &EntryPointUserOperationEventIterator{contract: _EntryPoint.contract, event: "UserOperationEvent", logs: logs, sub: sub}, nil
}

func (_EntryPoint *EntryPointFilterer) WatchUserOperationEvent(opts *bind.WatchOpts, sink chan<- *EntryPointUserOperationEvent, userOpHash [][32]byte, sender []common.Address, paymaster []common.Address) (event.Subscription, error) {
	var userOpHashRule []interface{}
	for _, userOpHashItem := range userOpHash {
		userOpHashRule = append(userOpHashRule, userOpHashItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var paymasterRule []interface{}
	for _, paymasterItem := range paymaster {
		paymasterRule = append(paymasterRule, paymasterItem)
	}

	logs, sub, err := _EntryPoint.contract.WatchLogs(opts, "UserOperationEvent", userOpHashRule, senderRule, paymasterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EntryPointUserOperationEvent)
				if err := _EntryPoint.contract.UnpackLog(event, "UserOperationEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EntryPoint *EntryPointFilterer) ParseUserOperationEvent(log types.Log) (*EntryPointUserOperationEvent, error) {
	event := new(EntryPointUserOperationEvent)
	if err := _EntryPoint.contract.UnpackLog(event, "UserOperationEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_EntryPoint *EntryPointFilterer) FilterUserOperationRevertReason(opts *bind.FilterOpts, userOpHash [][32]byte, sender []common.Address) (*EntryPointUserOperationRevertReasonIterator, error) {
	var userOpHashRule []interface{}
	for _, userOpHashItem := range userOpHash {
		userOpHashRule = append(userOpHashRule, userOpHashItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EntryPoint.contract.FilterLogs(opts, "UserOperationRevertReason", userOpHashRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EntryPointUserOperationRevertReasonIterator{contract: _EntryPoint.contract, event: "UserOperationRevertReason", logs: logs, sub: sub}, nil
}

func (_EntryPoint *EntryPointFilterer) WatchUserOperationRevertReason(opts *bind.WatchOpts, sink chan<- *EntryPointUserOperationRevertReason, userOpHash [][32]byte, sender []common.Address) (event.Subscription, error) {
	var userOpHashRule []interface{}
	for _, userOpHashItem := range userOpHash {
		userOpHashRule = append(userOpHashRule, userOpHashItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EntryPoint.contract.WatchLogs(opts, "UserOperationRevertReason", userOpHashRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EntryPointUserOperationRevertReason)
				if err := _EntryPoint.contract.UnpackLog(event, "UserOperationRevertReason", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EntryPoint *EntryPointFilterer) ParseUserOperationRevertReason(log types.Log) (*EntryPointUserOperationRevertReason, error) {
	event := new(EntryPointUserOperationRevertReason)
	if err := _EntryPoint.contract.UnpackLog(event, "UserOperationRevertReason", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_EntryPoint *EntryPointFilterer) FilterWithdrawn(opts *bind.FilterOpts, account []common.Address) (*EntryPointWithdrawnIterator, error) {
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EntryPoint.contract.FilterLogs(opts, "Withdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return &EntryPointWithdrawnIterator{contract: _EntryPoint.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

func (_EntryPoint *EntryPointFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *EntryPointWithdrawn, account []common.Address) (event.Subscription, error) {
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EntryPoint.contract.WatchLogs(opts, "Withdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EntryPointWithdrawn)
				if err := _EntryPoint.contract.UnpackLog(event, "Withdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EntryPoint *EntryPointFilterer) ParseWithdrawn(log types.Log) (*EntryPointWithdrawn, error) {
	event := new(EntryPointWithdrawn)
	if err := _EntryPoint.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_EntryPoint *EntryPoint) FilterLogs(opts *bind.FilterOpts, query ...[]interface{}) (*EntryPointEventIterator, error) {
	// Don't crash on a lazy user
	if opts == nil {
		opts = new(bind.FilterOpts)
	}
	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return nil, err
	}
	// Start the background filtering
	logs := make(chan types.Log, 128)

	config := ethereum.FilterQuery{
		Addresses: []common.Address{_EntryPoint.address},
		Topics:    topics,
		FromBlock: new(big.Int).SetUint64(opts.Start),
	}
	if opts.End != nil {
		config.ToBlock = new(big.Int).SetUint64(*opts.End)
	}
	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}
	buff, err := _EntryPoint.backend.FilterLogs(ctx, config)
	if err != nil {
		return nil, err
	}
	sub, err := event.NewSubscription(func(quit <-chan struct{}) error {
		for _, log := range buff {
			select {
			case logs <- log:
			case <-quit:
				return nil
			}
		}
		return nil
	}), nil
	if err != nil {
		return nil, err
	}

	return &EntryPointEventIterator{abi: _EntryPoint.abi, contract: _EntryPoint.contract, logs: logs, sub: sub}, nil
}

type EntryPointDecoder struct {
	abi abi.ABI
}

type DataError interface {
	Error() string
	ErrorData() interface{}
}

func (_EntryPoint *EntryPointDecoder) DecodeErrorFailedOp(orgErr error) (FailedOp, error) {
	var ret FailedOp

	if orgErr == nil {
		return ret, fmt.Errorf("cannot decode nil error")
	}

	de, ok := orgErr.(DataError)
	if !ok {
		return ret, fmt.Errorf("cannot read data from error: %s", orgErr)
	}

	dataHex, ok := de.ErrorData().(string)
	if !ok {
		return ret, fmt.Errorf("expect data type string got %T", de.ErrorData())
	}

	data, err := hexutil.Decode(dataHex)
	if err != nil {
		return ret, err
	}

	e := _EntryPoint.abi.Errors["FailedOp"]
	if len(data) < 4 {
		return ret, errors.New("invalid data for unpacking")
	}
	if !bytes.Equal(data[:4], e.ID[:4]) {
		return ret, errors.New("invalid data for unpacking")
	}
	unpacked, err := e.Inputs.Unpack(data[4:])
	if err != nil {
		return ret, err
	}
	err = e.Inputs.Copy(&ret, unpacked)
	return ret, err
}

func (_EntryPoint *EntryPointDecoder) DecodeErrorValidationResult(orgErr error) (ValidationResult, error) {
	var ret ValidationResult

	if orgErr == nil {
		return ret, fmt.Errorf("cannot decode nil error")
	}

	de, ok := orgErr.(DataError)
	if !ok {
		return ret, fmt.Errorf("cannot read data from error: %s", orgErr)
	}

	dataHex, ok := de.ErrorData().(string)
	if !ok {
		return ret, fmt.Errorf("expect data type string got %T", de.ErrorData())
	}

	data, err := hexutil.Decode(dataHex)
	if err != nil {
		return ret, err
	}

	e := _EntryPoint.abi.Errors["ValidationResult"]
	if len(data) < 4 {
		return ret, errors.New("invalid data for unpacking")
	}
	if !bytes.Equal(data[:4], e.ID[:4]) {
		return ret, errors.New("invalid data for unpacking")
	}
	unpacked, err := e.Inputs.Unpack(data[4:])
	if err != nil {
		return ret, err
	}
	fmt.Println(unpacked)
	err = e.Inputs.Copy(&ret, unpacked)
	return ret, err
}

func (_EntryPoint *EntryPointDecoder) DecodeErrorValidationResultWithAggregation(orgErr error) (ValidationResultWithAggregator, error) {
	var ret ValidationResultWithAggregator

	if orgErr == nil {
		return ret, fmt.Errorf("cannot decode nil error")
	}

	de, ok := orgErr.(DataError)
	if !ok {
		return ret, fmt.Errorf("cannot read data from error: %s", orgErr)
	}

	dataHex, ok := de.ErrorData().(string)
	if !ok {
		return ret, fmt.Errorf("expect data type string got %T", de.ErrorData())
	}

	data, err := hexutil.Decode(dataHex)
	if err != nil {
		return ret, err
	}

	e := _EntryPoint.abi.Errors["ValidationResultWithAggregation"]
	if len(data) < 4 {
		return ret, errors.New("invalid data for unpacking")
	}
	if !bytes.Equal(data[:4], e.ID[:4]) {
		return ret, errors.New("invalid data for unpacking")
	}
	unpacked, err := e.Inputs.Unpack(data[4:])
	if err != nil {
		return ret, err
	}
	err = e.Inputs.Copy(&ret, unpacked)
	return ret, err
}

func (_EntryPoint *EntryPointDecoder) DecodeErrorString(orgErr error) (string, error) {
	var ret string

	if orgErr == nil {
		return ret, fmt.Errorf("cannot decode nil error")
	}

	de, ok := orgErr.(DataError)
	if !ok {
		return ret, fmt.Errorf("cannot read data from error: %s", orgErr)
	}

	dataHex, ok := de.ErrorData().(string)
	if !ok {
		return ret, fmt.Errorf("expect data type string got %T", de.ErrorData())
	}

	data, err := hexutil.Decode(dataHex)
	if err != nil {
		return ret, err
	}

	return abi.UnpackRevert(data)
}

func (_EntryPoint *EntryPointDecoder) DecodeHandleOps(data []byte) ([]UserOperation, common.Address, error) {
	method := _EntryPoint.abi.Methods["handleOps"]

	if !bytes.Equal(data[:4], method.ID) {
		return nil, common.Address{}, errors.New("invalid data for unpacking")
	}

	res, err := method.Inputs.Unpack(data[4:])
	if err != nil {
		return nil, common.Address{}, err
	}

	if len(res) != 2 {
		return nil, common.Address{}, errors.New("invalid data for unpacking")
	}

	beneficiary, ok := abi.ConvertType(res[1], new(common.Address)).(*common.Address)
	if !ok {
		return nil, common.Address{}, errors.New("cannot convert type to common.Address")
	}

	fmt.Printf("%T", res[0])

	var x []UserOperation

	fmt.Printf("%T", x)

	userOps, ok := abi.ConvertType(res[0], new([]UserOperation)).(*[]UserOperation)
	if !ok {
		return nil, common.Address{}, errors.New("cannot convert type to []UserOperation")
	}

	return *userOps, *beneficiary, nil
}

var _ IEntryPoint = (*EntryPoint)(nil)

type EntryPoint struct {
	address    common.Address
	abi        abi.ABI
	backend    bind.ContractBackend
	contract   *bind.BoundContract
	caller     IEntryPointCaller
	transactor IEntryPointTransactor
	filterer   IEntryPointFilterer
	decoder    IEntryPointDecoder
}

func (_EntryPoint *EntryPoint) Address() common.Address {
	return _EntryPoint.address
}

func (_EntryPoint *EntryPoint) Caller() IEntryPointCaller {
	return _EntryPoint.caller
}

func (_EntryPoint *EntryPoint) Transactor() IEntryPointTransactor {
	return _EntryPoint.transactor
}

func (_EntryPoint *EntryPoint) Filterer() IEntryPointFilterer {
	return _EntryPoint.filterer
}

func (_EntryPoint *EntryPoint) Decoder() IEntryPointDecoder {
	return _EntryPoint.decoder
}

func NewEntryPoint(address common.Address, backend bind.ContractBackend) (*EntryPoint, error) {
	abi, err := abi.JSON(strings.NewReader(EntryPointABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindEntryPoint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EntryPoint{
		address:    address,
		abi:        abi,
		contract:   contract,
		backend:    backend,
		caller:     &EntryPointCaller{contract: contract},
		transactor: &EntryPointTransactor{contract: contract},
		filterer:   &EntryPointFilterer{contract: contract},
		decoder:    &EntryPointDecoder{abi: abi},
	}, nil
}

func NewEntryPointCaller(address common.Address, caller bind.ContractCaller) (*EntryPointCaller, error) {
	contract, err := bindEntryPoint(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EntryPointCaller{contract: contract}, nil
}

func NewEntryPointTransactor(address common.Address, transactor bind.ContractTransactor) (*EntryPointTransactor, error) {
	contract, err := bindEntryPoint(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EntryPointTransactor{contract: contract}, nil
}

func NewEntryPointFilterer(address common.Address, filterer bind.ContractFilterer) (*EntryPointFilterer, error) {
	contract, err := bindEntryPoint(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EntryPointFilterer{contract: contract}, nil
}

func bindEntryPoint(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EntryPointABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}
