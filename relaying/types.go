package relaying

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type RelayTransactionArgs struct {
	Wallet        *common.Address `json:"wallet,omitempty"`
	Data          *hexutil.Bytes  `json:"data,omitempty"`
	Nonce         *hexutil.Big    `json:"nonce,omitempty"`
	Signatures    *hexutil.Bytes  `json:"signatures,omitempty"`
	GasPrice      *hexutil.Big    `json:"gasPrice,omitempty"`
	GasLimit      *hexutil.Big    `json:"gasLimit,omitempty"`
	RefundToken   *common.Address `json:"refundToken,omitempty"`
	RefundAddress *common.Address `json:"refundAddress,omitempty"`
}

type RelayTransaction struct {
	Wallet        common.Address
	Data          []byte
	Nonce         *big.Int
	Signatures    []byte
	GasPrice      *big.Int
	GasLimit      *big.Int
	RefundToken   common.Address
	RefundAddress common.Address
}
