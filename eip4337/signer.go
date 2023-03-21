package eip4337

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/eip4337/log"
	"github.com/evmos/ethermint/rpc/backend"
)

type ISigner interface {
	GetAddress() common.Address
	SignTx(tx *ethtypes.Transaction) (*ethtypes.Transaction, error)
}

var _ ISigner = (*Signer)(nil)

type Signer struct {
	logger  log.Logger
	keyring keyring.Keyring
	address common.Address
	backend backend.EVMBackend
}

func NewSigner(logger log.Logger, keyring keyring.Keyring, address common.Address, backend backend.EVMBackend) ISigner {
	return &Signer{
		logger:  log.EnsureLogger(logger),
		keyring: keyring,
		address: address,
		backend: backend,
	}
}

func (signer *Signer) GetAddress() common.Address {
	return signer.address
}

func (signer *Signer) SignTx(tx *ethtypes.Transaction) (*ethtypes.Transaction, error) {
	bn, err := signer.backend.BlockNumber()
	if err != nil {
		return nil, err
	}

	s := ethtypes.MakeSigner(signer.backend.ChainConfig(), new(big.Int).SetUint64(uint64(bn)))
	txHash := s.Hash(tx)

	sig, _, err := signer.keyring.SignByAddress(sdk.AccAddress(signer.address.Bytes()), txHash.Bytes())
	if err != nil {
		return nil, err
	}

	tx, err = tx.WithSignature(s, sig)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
