package relaying

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type RelayerConfig struct {
	Address                common.Address
	SenderAddresses        []common.Address
	AllowedRefundAddresses []common.Address
	AllowedRefundTokens    []common.Address
	MinGasPrices           []*big.Int
	GasMultiplier          float64
}

func (config RelayerConfig) Validate() error {
	if len(config.SenderAddresses) == 0 {
		return errors.New("empty sender list")
	}
	if len(config.AllowedRefundAddresses) == 0 {
		return errors.New("empty refund address list")
	}
	if len(config.AllowedRefundTokens) == 0 {
		return errors.New("empty refund token list")
	}
	if len(config.MinGasPrices) != len(config.AllowedRefundTokens) {
		return errors.New("refund tokens and min gas prices length not match")
	}

	if hasDuplicate(config.SenderAddresses) {
		return errors.New("sender list contains duplicate")
	}
	if hasDuplicate(config.AllowedRefundAddresses) {
		return errors.New("refund address list contains duplicate")
	}
	if hasDuplicate(config.AllowedRefundTokens) {
		return errors.New("refund token list contains duplicate")
	}

	for _, minGasPrice := range config.MinGasPrices {
		if minGasPrice == nil || minGasPrice.Cmp(new(big.Int).SetInt64(0)) < 0 {
			return errors.New("minGasPrice cannot be negative")
		}
	}
	if config.GasMultiplier <= 0 {
		return errors.New("gasMultiplier must be positive")
	}
	return nil
}

func hasDuplicate(arr []common.Address) bool {
	m := make(map[string]struct{})
	for _, a := range arr {
		if _, ok := m[a.Hex()]; ok {
			return true
		}
		m[a.Hex()] = struct{}{}
	}
	return false
}
