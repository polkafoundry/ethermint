package relaying

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/rpc/backend"
)

type IKeyStore interface {
	SignTx(fromAddress common.Address, tx *ethtypes.Transaction) (*ethtypes.Transaction, error)
}

type KeyStore struct {
	keyring keyring.Keyring
	backend backend.EVMBackend
}

func NewKeyStore(keyring keyring.Keyring, backend backend.EVMBackend) IKeyStore {
	return &KeyStore{
		keyring: keyring,
		backend: backend,
	}
}

func (ks *KeyStore) SignTx(fromAddress common.Address, tx *ethtypes.Transaction) (*ethtypes.Transaction, error) {
	bn, err := ks.backend.BlockNumber()
	if err != nil {
		return nil, err
	}

	s := ethtypes.MakeSigner(ks.backend.ChainConfig(), new(big.Int).SetUint64(uint64(bn)))
	txHash := s.Hash(tx)

	sig, _, err := ks.keyring.SignByAddress(sdk.AccAddress(fromAddress.Bytes()), txHash.Bytes())
	if err != nil {
		return nil, err
	}

	tx, err = tx.WithSignature(s, sig)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
