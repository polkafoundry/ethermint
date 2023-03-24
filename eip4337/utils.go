package eip4337

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/eip4337/types"
	rpctypes "github.com/evmos/ethermint/rpc/types"
	"golang.org/x/crypto/sha3"
)

// hasherPool holds LegacyKeccak256 hashers.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

func GetUserOpHashes(userOps []types.UserOperation, entryPointAddress common.Address, chainID *big.Int) []common.Hash {
	var hashes []common.Hash
	for _, op := range userOps {
		hashes = append(hashes, GetUserOpHash(op, entryPointAddress, chainID))
	}
	return hashes
}

func GetUserOpHash(userOp types.UserOperation, entryPointAddress common.Address, chainID *big.Int) common.Hash {
	var h common.Hash
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)

	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)

	args := abi.Arguments{
		{
			Type: bytes32Type,
		},
		{
			Type: addressType,
		},
		{
			Type: uint256Type,
		},
	}

	var innerHash common.Hash
	sha.Reset()
	packed := PackUserOp(userOp, true)
	_, _ = sha.Write(packed)
	_, _ = sha.Read(innerHash[:])

	sha.Reset()
	bz, _ := args.Pack(innerHash, entryPointAddress, chainID)
	_, _ = sha.Write(bz)
	_, _ = sha.Read(h[:])
	return h
}

func PackUserOp(op types.UserOperation, forSignature bool) []byte {
	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)
	bytesType, _ := abi.NewType("bytes", "", nil)

	arguments := abi.Arguments{
		{Type: addressType}, // sender: *common.Address
		{Type: uint256Type}, // nonce: *big.Int
		{Type: bytesType},   // initCode: []byte
		{Type: bytesType},   // callData: []byte
		{Type: uint256Type}, // callGasLimit: *big.Int
		{Type: uint256Type}, // verificationGasLimit: *big.Int
		{Type: uint256Type}, // preVerificationGas: *bit.Int
		{Type: uint256Type}, // maxFeePerGas: *big.Int
		{Type: uint256Type}, // maxPriorityFeePerGas: *big.Int
		{Type: bytesType},   // paymasterAndData: []byte
		{Type: bytesType},   // signature: []byte
	}

	values := []interface{}{
		op.Sender(),
		op.Nonce(),
		op.InitCode(),
		op.CallData(),
		op.CallGasLimit(),
		op.VerificationGasLimit(),
		op.PreVerificationGas(),
		op.MaxFeePerGas(),
		op.MaxPriorityFeePerGas(),
		op.PaymasterAndData(),
	}

	if forSignature {
		values = append(values, []byte{})
		bz, _ := arguments.Pack(values...)
		return bz[:len(bz)-32]
	}

	values = append(values, op.Signature())
	bz, _ := arguments.Pack(values...)

	return bz
}

func newDummyBytesSliceWithValue(n int, v byte) []byte {
	arr := make([]byte, n)
	for i := range arr {
		arr[i] = v
	}
	return arr
}

func toBlockNumberArg(bn *big.Int) rpctypes.BlockNumber {
	if bn == nil || !bn.IsInt64() {
		return rpctypes.EthLatestBlockNumber
	}
	return rpctypes.BlockNumber(bn.Int64())
}

// addrToAddrPtr return a pointer to a copy of an address.
// Zero address is considered nil
func addrToAddrPtr(addr common.Address) *common.Address {
	if addr.String() == "0x0000000000000000000000000000000000000000" {
		return nil
	}
	cpy := addr
	return &cpy
}

func toHexUtilUint64(v uint64) *hexutil.Uint64 {
	if v == 0 {
		return nil
	}
	cpy := v
	return (*hexutil.Uint64)(&cpy)
}
