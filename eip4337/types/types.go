package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
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
	Sender               *common.Address `json:"sender,omitempty"`
	Nonce                *hexutil.Big    `json:"nonce,omitempty"`
	InitCode             *hexutil.Bytes  `json:"initCode,omitempty"`
	CallData             *hexutil.Bytes  `json:"callData,omitempty"`
	CallGasLimit         *hexutil.Big    `json:"callGasLimit,omitempty"`
	VerificationGasLimit *hexutil.Big    `json:"verificationGasLimit,omitempty"`
	PreVerificationGas   *hexutil.Big    `json:"preVerificationGas,omitempty"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas,omitempty"`
	PaymasterAndData     *hexutil.Bytes  `json:"paymasterAndData,omitempty"`
	Signature            *hexutil.Bytes  `json:"signature,omitempty"`
}

func NewUserOperation(args UserOperationArgs) UserOperation {
	userOp := UserOperation{
		sender:               copyAddressPtr(args.Sender),
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
		nonce:                nil,
		initCode:             common.CopyBytes(op.initCode),
		callData:             common.CopyBytes(op.callData),
		callGasLimit:         nil,
		verificationGasLimit: nil,
		preVerificationGas:   nil,
		maxFeePerGas:         nil,
		maxPriorityFeePerGas: nil,
		paymasterAndData:     common.CopyBytes(op.paymasterAndData),
		signature:            common.CopyBytes(op.signature),
	}
	if op.nonce != nil {
		cpy.nonce = new(big.Int).Set(op.nonce)
	}
	if op.callGasLimit != nil {
		cpy.callGasLimit = new(big.Int).Set(op.callGasLimit)
	}
	if op.verificationGasLimit != nil {
		cpy.verificationGasLimit = new(big.Int).Set(op.verificationGasLimit)
	}
	if op.preVerificationGas != nil {
		cpy.preVerificationGas = new(big.Int).Set(op.preVerificationGas)
	}
	if op.maxFeePerGas != nil {
		cpy.maxFeePerGas = new(big.Int).Set(op.maxFeePerGas)
	}
	if op.maxPriorityFeePerGas != nil {
		cpy.maxPriorityFeePerGas = new(big.Int).Set(op.maxPriorityFeePerGas)
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

type EstimateUserOpGasResult struct {
	PreVerificationGas *hexutil.Big    `json:"preVerificationGas,omitempty"`
	VerificationGas    *hexutil.Big    `json:"verificationGas,omitempty"`
	ValidAfter         *hexutil.Uint64 `json:"validAfter,omitempty"`
	ValidUntil         *hexutil.Uint64 `json:"validUntil,omitempty"`
	CallGasLimit       *hexutil.Uint64 `json:"callGasLimit,omitempty"`
}

type UserOperationReceipt struct {
	UserOpHash    *common.Hash      `json:"userOpHash,omitempty"`
	Sender        *common.Address   `json:"sender,omitempty"`
	Nonce         *hexutil.Big      `json:"nonce,omitempty"`
	Paymaster     *common.Address   `json:"paymaster,omitempty"`
	ActualGasCode *hexutil.Big      `json:"actualGasCode,omitempty"`
	ActualGasUsed *hexutil.Big      `json:"actualGasUsed,omitempty"`
	Success       bool              `json:"success,omitempty"`
	Reason        *string           `json:"reason,omitempty"`
	Logs          []*ethtypes.Log   `json:"logs,omitempty"`
	Receipt       *ethtypes.Receipt `json:"receipt,omitempty"`
}

type UserOperationResponse struct {
	Sender               *common.Address `json:"sender,omitempty"`
	Nonce                *hexutil.Big    `json:"nonce,omitempty"`
	InitCode             *hexutil.Bytes  `json:"initCode,omitempty"`
	CallData             *hexutil.Bytes  `json:"callData,omitempty"`
	CallGasLimit         *hexutil.Big    `json:"callGasLimit,omitempty"`
	VerificationGasLimit *hexutil.Big    `json:"verificationGasLimit,omitempty"`
	PreVerificationGas   *hexutil.Big    `json:"preVerificationGas,omitempty"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas,omitempty"`
	PaymasterAndData     *hexutil.Bytes  `json:"paymasterAndData,omitempty"`
	Signature            *hexutil.Bytes  `json:"signature,omitempty"`
}

type UserOperationByHashResponse struct {
	UserOperation   UserOperationResponse `json:"userOperation"`
	EntryPoint      *common.Address       `json:"entryPoint"`
	BlockNumber     *hexutil.Big          `json:"blockNumber"`
	BlockHash       *common.Hash          `json:"blockHash"`
	TransactionHash *common.Hash          `json:"transactionHash"`
}

func ToUserOperationResponse(userOp UserOperation) UserOperationResponse {
	initCode := userOp.InitCode()
	callData := userOp.CallData()
	paymasterAndData := userOp.PaymasterAndData()
	signature := userOp.Signature()

	return UserOperationResponse{
		Sender:               userOp.Sender(),
		Nonce:                (*hexutil.Big)(userOp.Nonce()),
		InitCode:             (*hexutil.Bytes)(&initCode),
		CallData:             (*hexutil.Bytes)(&callData),
		CallGasLimit:         (*hexutil.Big)(userOp.CallGasLimit()),
		VerificationGasLimit: (*hexutil.Big)(userOp.VerificationGasLimit()),
		PreVerificationGas:   (*hexutil.Big)(userOp.PreVerificationGas()),
		MaxFeePerGas:         (*hexutil.Big)(userOp.MaxFeePerGas()),
		MaxPriorityFeePerGas: (*hexutil.Big)(userOp.MaxPriorityFeePerGas()),
		PaymasterAndData:     (*hexutil.Bytes)(&paymasterAndData),
		Signature:            (*hexutil.Bytes)(&signature),
	}
}

func ToUserOperationResponses(userOps []UserOperation) []UserOperationResponse {
	userOpResponses := make([]UserOperationResponse, len(userOps))
	for idx, op := range userOps {
		userOpResponses[idx] = ToUserOperationResponse(op)
	}
	return userOpResponses
}

type ReputationArgs struct {
	Address     common.Address `json:"address,omitempty"`
	OpsSeen     int64          `json:"opsSeen,omitempty"`
	OpsIncluded int64          `json:"opsIncluded,omitempty"`
}

type ReputationResponse struct {
	Address     common.Address `json:"address,omitempty"`
	OpsSeen     int64          `json:"opsSeen"`
	OpsIncluded int64          `json:"opsIncluded"`
	Status      string         `json:"status"`
}
