package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type UserOperation struct {
	sender               *common.Address
	nonce                *big.Int
	initCode             []byte
	callData             []byte
	callGasLimit         *big.Int
	verificationGasLimit *big.Int
	preVerificationGas   *big.Int
	maxFeePerGas         *big.Int
	maxPriorityFeePerGas *big.Int
	paymasterAndData     []byte
	signature            []byte
}

type UserOperationArgs struct {
	Sender               *common.Address
	Nonce                *hexutil.Big
	InitCode             *hexutil.Bytes
	CallData             *hexutil.Bytes
	CallGasLimit         *hexutil.Big
	VerificationGasLimit *hexutil.Big
	PreVerificationGas   *hexutil.Big
	MaxFeePerGas         *hexutil.Big
	MaxPriorityFeePerGas *hexutil.Big
	PaymasterAndData     *hexutil.Bytes
	Signature            *hexutil.Bytes
}

func NewUserOperation(args UserOperationArgs) UserOperation {
	userOp := UserOperation{
		sender:               args.Sender,
		nonce:                (*big.Int)(args.Nonce),
		initCode:             nil,
		callData:             nil,
		callGasLimit:         (*big.Int)(args.CallGasLimit),
		verificationGasLimit: (*big.Int)(args.VerificationGasLimit),
		preVerificationGas:   (*big.Int)(args.PreVerificationGas),
		maxFeePerGas:         (*big.Int)(args.MaxFeePerGas),
		maxPriorityFeePerGas: (*big.Int)(args.MaxPriorityFeePerGas),
		paymasterAndData:     nil,
		signature:            nil,
	}

	if args.InitCode != nil {
		userOp.initCode = common.CopyBytes(*args.InitCode)
	}

	if args.CallData != nil {
		userOp.callData = common.CopyBytes(*args.CallData)
	}

	if args.PaymasterAndData != nil {
		userOp.paymasterAndData = common.CopyBytes(*args.PaymasterAndData)
	}

	if args.Signature != nil {
		userOp.signature = common.CopyBytes(*args.Signature)
	}

	return userOp
}

func (op UserOperation) WithSignature(signature []byte) UserOperation {
	cpy := op.copy()
	cpy.signature = common.CopyBytes(signature)
	return cpy
}

func (op UserOperation) WithPreVerificationGas(preVerificationGas *big.Int) UserOperation {
	cpy := op.copy()
	cpy.preVerificationGas = nil
	if preVerificationGas != nil {
		cpy.preVerificationGas = new(big.Int).Set(preVerificationGas)
	}
	return cpy
}

func (op UserOperation) Sender() *common.Address { return copyAddressPtr(op.sender) }

func (op UserOperation) Nonce() *big.Int { return new(big.Int).Set(op.nonce) }

func (op UserOperation) InitCode() []byte { return common.CopyBytes(op.initCode) }

func (op UserOperation) CallData() []byte { return common.CopyBytes(op.callData) }

func (op UserOperation) CallGasLimit() *big.Int {
	return new(big.Int).Set(op.callGasLimit)
}

func (op UserOperation) VerificationGasLimit() *big.Int {
	return new(big.Int).Set(op.verificationGasLimit)
}

func (op UserOperation) PreVerificationGas() *big.Int {
	return new(big.Int).Set(op.preVerificationGas)
}

func (op UserOperation) MaxFeePerGas() *big.Int {
	return new(big.Int).Set(op.maxFeePerGas)
}

func (op UserOperation) MaxPriorityFeePerGas() *big.Int {
	return new(big.Int).Set(op.maxPriorityFeePerGas)
}

func (op UserOperation) PaymasterAndData() []byte { return common.CopyBytes(op.paymasterAndData) }

func (op UserOperation) Signature() []byte { return common.CopyBytes(op.signature) }

func (op UserOperation) copy() UserOperation {
	cpy := UserOperation{
		sender:               copyAddressPtr(op.sender),
		nonce:                new(big.Int),
		initCode:             common.CopyBytes(op.initCode),
		callData:             common.CopyBytes(op.callData),
		callGasLimit:         new(big.Int),
		verificationGasLimit: new(big.Int),
		preVerificationGas:   new(big.Int),
		maxFeePerGas:         new(big.Int),
		maxPriorityFeePerGas: new(big.Int),
		paymasterAndData:     common.CopyBytes(op.paymasterAndData),
		signature:            common.CopyBytes(op.signature),
	}
	if op.nonce != nil {
		cpy.nonce.Set(op.nonce)
	}
	if op.callGasLimit != nil {
		cpy.callGasLimit.Set(op.callGasLimit)
	}
	if op.verificationGasLimit != nil {
		cpy.verificationGasLimit.Set(op.verificationGasLimit)
	}
	if op.preVerificationGas != nil {
		cpy.preVerificationGas.Set(op.preVerificationGas)
	}
	if op.maxFeePerGas != nil {
		cpy.maxFeePerGas.Set(op.maxFeePerGas)
	}
	if op.maxPriorityFeePerGas != nil {
		cpy.maxPriorityFeePerGas.Set(op.maxPriorityFeePerGas)
	}
	return cpy
}

// copyAddressPtr copies an address.
func copyAddressPtr(a *common.Address) *common.Address {
	if a == nil {
		return nil
	}
	cpy := *a
	return &cpy
}

type ReferencedCodeHashes struct {
	Addresses []common.Address
	Hash      common.Hash
}
