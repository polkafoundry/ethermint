package encoding

import (
	"github.com/cosmos/cosmos-sdk/client"
	amino "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	enccodec "github.com/evmos/ethermint/encoding/codec"
)

// MakeConfig creates an EncodingConfig for testing
func MakeConfig(mb module.BasicManager) params.EncodingConfig {
	cdc := amino.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	codec := amino.NewProtoCodec(interfaceRegistry)

	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          tx.NewTxConfig(codec, tx.DefaultSignModes),
		Amino:             cdc,
	}

	enccodec.RegisterLegacyAminoCodec(encodingConfig.Amino)
	mb.RegisterLegacyAminoCodec(encodingConfig.Amino)
	enccodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	mb.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

type OpConfig struct {
	OpEncodingConfig client.OpEncodingConfig
}

type opEncodingConfig struct {
	opEncoder sdk.OpEncoder
	opDecoder sdk.OpDecoder
}

func (ec opEncodingConfig) OpEncoder() sdk.OpEncoder {
	return ec.opEncoder
}

func (ec opEncodingConfig) OpDecoder() sdk.OpDecoder {
	return ec.opDecoder
}

func MakeOpConfig() OpConfig {
	return OpConfig{
		OpEncodingConfig: &opEncodingConfig{
			opEncoder: sdk.DefaultOpEncoder(),
			opDecoder: sdk.DefaultOpDecoder(),
		},
	}
}
