package relaying

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	rpctypes "github.com/evmos/ethermint/rpc/types"
)

func ensureBytes(bz *hexutil.Bytes) []byte {
	if bz == nil {
		return nil
	}
	return *bz
}

func ensureAddress(addr *common.Address) common.Address {
	if addr == nil {
		return common.Address{}
	}
	return *addr
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

func addressInArray(arr []common.Address, addr common.Address) bool {
	for _, a := range arr {
		if a == addr {
			return true
		}
	}
	return false
}
