package types

import (
	"bytes"
	_ "embed"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	//go:embed EntryPointContractABI.json
	entryPointABIJSON []byte

	EntryPointContractABI abi.ABI
)

func init() {
	var err error
	EntryPointContractABI, err = abi.JSON(bytes.NewReader(entryPointABIJSON))
	if err != nil {
		panic(err)
	}
}
