package utils

import "github.com/ethereum/go-ethereum/common"

// CopyAddressPtr copies an address.
func CopyAddressPtr(a *common.Address) *common.Address {
	if a == nil {
		return nil
	}
	cpy := *a
	return &cpy
}
