package config

import (
	"errors"
	"fmt"
	"math/big"
	"path"
	"time"

	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/server/config"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// DefaultGRPCAddress is the default address the gRPC server binds to.
	DefaultGRPCAddress = "0.0.0.0:9900"

	// DefaultJSONRPCAddress is the default address the JSON-RPC server binds to.
	DefaultJSONRPCAddress = "0.0.0.0:8545"

	// DefaultJSONRPCWsAddress is the default address the JSON-RPC WebSocket server binds to.
	DefaultJSONRPCWsAddress = "0.0.0.0:8546"

	// DefaultJsonRPCMetricsAddress is the default address the JSON-RPC Metrics server binds to.
	DefaultJSONRPCMetricsAddress = "0.0.0.0:6065"

	// DefaultEVMTracer is the default vm.Tracer type
	DefaultEVMTracer = ""

	DefaultMaxTxGasWanted = 0

	DefaultGasCap uint64 = 25000000

	DefaultFilterCap int32 = 200

	DefaultFeeHistoryCap int32 = 100

	DefaultLogsCap int32 = 10000

	DefaultBlockRangeCap int32 = 10000

	DefaultEVMTimeout = 5 * time.Second
	// default 1.0 eth
	DefaultTxFeeCap float64 = 1.0

	DefaultHTTPTimeout = 30 * time.Second

	DefaultHTTPIdleTimeout = 120 * time.Second
	// DefaultAllowUnprotectedTxs value is false
	DefaultAllowUnprotectedTxs = false
	// DefaultMaxOpenConnections represents the amount of open connections (unlimited = 0)
	DefaultMaxOpenConnections = 0
)

var evmTracers = []string{"json", "markdown", "struct", "access_list"}

// Config defines the server's top level configuration. It includes the default app config
// from the SDK as well as the EVM configuration to enable the JSON-RPC APIs.
type Config struct {
	config.Config

	EVM     EVMConfig     `mapstructure:"evm"`
	JSONRPC JSONRPCConfig `mapstructure:"json-rpc"`
	TLS     TLSConfig     `mapstructure:"tls"`
	Bundler BundlerConfig `mapstructure:"bundler"`
}

type BundlerConfig struct {
	Enable             bool   `mapstructure:"enable"`
	Debug              bool   `mapstructure:"debug"`
	BeneficiaryAddress string `mapstructure:"beneficiary"`
	EntryPointAddress  string `mapstructure:"entrypoint"`
	MinBalance         string `mapstructure:"min-balance"`
	SignerAddress      string `mapstructure:"signer"`

	AutoBundle            bool     `mapstructure:"auto-bundle"`
	Whitelist             []string `mapstructure:"whitelist"`
	Blacklist             []string `mapstructure:"blacklist"`
	MaxBundleGas          uint64   `mapstructure:"max-bundle-gas"`
	MinStake              string   `mapstructure:"min-stake"`
	MinUnstakeDelay       uint64   `mapstructure:"min-unstake-delay"`
	AutoBundleInterval    uint64   `mapstructure:"auto-bundle-interval"`
	AutoBundleMempoolSize uint64   `mapstructure:"auto-bundle-mempool-size"`
}

func DefaultBundlerConfig() *BundlerConfig {
	return &BundlerConfig{
		Enable:                false,
		Debug:                 false,
		BeneficiaryAddress:    "",
		EntryPointAddress:     "",
		MinBalance:            "0x",
		SignerAddress:         "",
		AutoBundle:            false,
		Whitelist:             []string{"0xcafebabe"},
		Blacklist:             []string{"deadbeef"},
		MaxBundleGas:          0,
		MinStake:              "0x",
		MinUnstakeDelay:       0,
		AutoBundleInterval:    0,
		AutoBundleMempoolSize: 0,
	}
}

func (config *BundlerConfig) Validate() error {
	// do not check other config values if bundler is not enabled
	if !config.Enable {
		return nil
	}

	if len(config.MinBalance) < 2 {
		return errors.New("min-balance must be a number in a hex string format")
	}
	_, ok := new(big.Int).SetString(config.MinBalance[2:], 16)
	if !ok {
		return errors.New("min-balance must be a number in a hex string format")
	}

	if len(config.MinStake) < 2 {
		return errors.New("min-stake must be a number in a hex string format")
	}
	_, ok = new(big.Int).SetString(config.MinStake[2:], 16)
	if !ok {
		return errors.New("min-stake must be a number in a hex string format")
	}

	// FIXME: add more validation
	return nil
}

// EVMConfig defines the application configuration values for the EVM.
type EVMConfig struct {
	// Tracer defines vm.Tracer type that the EVM will use if the node is run in
	// trace mode. Default: 'json'.
	Tracer string `mapstructure:"tracer"`
	// MaxTxGasWanted defines the gas wanted for each eth tx returned in ante handler in check tx mode.
	MaxTxGasWanted uint64 `mapstructure:"max-tx-gas-wanted"`
}

// JSONRPCConfig defines configuration for the EVM RPC server.
type JSONRPCConfig struct {
	// API defines a list of JSON-RPC namespaces that should be enabled
	API []string `mapstructure:"api"`
	// Address defines the HTTP server to listen on
	Address string `mapstructure:"address"`
	// WsAddress defines the WebSocket server to listen on
	WsAddress string `mapstructure:"ws-address"`
	// GasCap is the global gas cap for eth-call variants.
	GasCap uint64 `mapstructure:"gas-cap"`
	// EVMTimeout is the global timeout for eth-call.
	EVMTimeout time.Duration `mapstructure:"evm-timeout"`
	// TxFeeCap is the global tx-fee cap for send transaction
	TxFeeCap float64 `mapstructure:"txfee-cap"`
	// FilterCap is the global cap for total number of filters that can be created.
	FilterCap int32 `mapstructure:"filter-cap"`
	// FeeHistoryCap is the global cap for total number of blocks that can be fetched
	FeeHistoryCap int32 `mapstructure:"feehistory-cap"`
	// Enable defines if the EVM RPC server should be enabled.
	Enable bool `mapstructure:"enable"`
	// LogsCap defines the max number of results can be returned from single `eth_getLogs` query.
	LogsCap int32 `mapstructure:"logs-cap"`
	// BlockRangeCap defines the max block range allowed for `eth_getLogs` query.
	BlockRangeCap int32 `mapstructure:"block-range-cap"`
	// HTTPTimeout is the read/write timeout of http json-rpc server.
	HTTPTimeout time.Duration `mapstructure:"http-timeout"`
	// HTTPIdleTimeout is the idle timeout of http json-rpc server.
	HTTPIdleTimeout time.Duration `mapstructure:"http-idle-timeout"`
	// AllowUnprotectedTxs restricts unprotected (non EIP155 signed) transactions to be submitted via
	// the node's RPC when global parameter is disabled.
	AllowUnprotectedTxs bool `mapstructure:"allow-unprotected-txs"`
	// MaxOpenConnections sets the maximum number of simultaneous connections
	// for the server listener.
	MaxOpenConnections int `mapstructure:"max-open-connections"`
	// EnableIndexer defines if enable the custom indexer service.
	EnableIndexer bool `mapstructure:"enable-indexer"`
	// MetricsAddress defines the metrics server to listen on
	MetricsAddress string `mapstructure:"metrics-address"`
}

// TLSConfig defines the certificate and matching private key for the server.
type TLSConfig struct {
	// CertificatePath the file path for the certificate .pem file
	CertificatePath string `mapstructure:"certificate-path"`
	// KeyPath the file path for the key .pem file
	KeyPath string `mapstructure:"key-path"`
}

// AppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func AppConfig(denom string) (string, interface{}) {
	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := config.DefaultConfig()

	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In ethermint, we set the min gas prices to 0.
	if denom != "" {
		srvCfg.MinGasPrices = "0" + denom
	}

	customAppConfig := Config{
		Config:  *srvCfg,
		EVM:     *DefaultEVMConfig(),
		JSONRPC: *DefaultJSONRPCConfig(),
		TLS:     *DefaultTLSConfig(),
		Bundler: *DefaultBundlerConfig(),
	}

	customAppTemplate := config.DefaultConfigTemplate + DefaultConfigTemplate

	return customAppTemplate, customAppConfig
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{
		Config:  *config.DefaultConfig(),
		EVM:     *DefaultEVMConfig(),
		JSONRPC: *DefaultJSONRPCConfig(),
		TLS:     *DefaultTLSConfig(),
		Bundler: *DefaultBundlerConfig(),
	}
}

// DefaultEVMConfig returns the default EVM configuration
func DefaultEVMConfig() *EVMConfig {
	return &EVMConfig{
		Tracer:         DefaultEVMTracer,
		MaxTxGasWanted: DefaultMaxTxGasWanted,
	}
}

// Validate returns an error if the tracer type is invalid.
func (c EVMConfig) Validate() error {
	if c.Tracer != "" && !strings.StringInSlice(c.Tracer, evmTracers) {
		return fmt.Errorf("invalid tracer type %s, available types: %v", c.Tracer, evmTracers)
	}

	return nil
}

// GetDefaultAPINamespaces returns the default list of JSON-RPC namespaces that should be enabled
func GetDefaultAPINamespaces() []string {
	return []string{"eth", "net", "web3"}
}

// GetAPINamespaces returns the all the available JSON-RPC API namespaces.
func GetAPINamespaces() []string {
	return []string{"web3", "eth", "personal", "net", "txpool", "debug", "miner"}
}

// DefaultJSONRPCConfig returns an EVM config with the JSON-RPC API enabled by default
func DefaultJSONRPCConfig() *JSONRPCConfig {
	return &JSONRPCConfig{
		Enable:              true,
		API:                 GetDefaultAPINamespaces(),
		Address:             DefaultJSONRPCAddress,
		WsAddress:           DefaultJSONRPCWsAddress,
		GasCap:              DefaultGasCap,
		EVMTimeout:          DefaultEVMTimeout,
		TxFeeCap:            DefaultTxFeeCap,
		FilterCap:           DefaultFilterCap,
		FeeHistoryCap:       DefaultFeeHistoryCap,
		BlockRangeCap:       DefaultBlockRangeCap,
		LogsCap:             DefaultLogsCap,
		HTTPTimeout:         DefaultHTTPTimeout,
		HTTPIdleTimeout:     DefaultHTTPIdleTimeout,
		AllowUnprotectedTxs: DefaultAllowUnprotectedTxs,
		MaxOpenConnections:  DefaultMaxOpenConnections,
		EnableIndexer:       false,
		MetricsAddress:      DefaultJSONRPCMetricsAddress,
	}
}

// Validate returns an error if the JSON-RPC configuration fields are invalid.
func (c JSONRPCConfig) Validate() error {
	if c.Enable && len(c.API) == 0 {
		return errors.New("cannot enable JSON-RPC without defining any API namespace")
	}

	if c.FilterCap < 0 {
		return errors.New("JSON-RPC filter-cap cannot be negative")
	}

	if c.FeeHistoryCap <= 0 {
		return errors.New("JSON-RPC feehistory-cap cannot be negative or 0")
	}

	if c.TxFeeCap < 0 {
		return errors.New("JSON-RPC tx fee cap cannot be negative")
	}

	if c.EVMTimeout < 0 {
		return errors.New("JSON-RPC EVM timeout duration cannot be negative")
	}

	if c.LogsCap < 0 {
		return errors.New("JSON-RPC logs cap cannot be negative")
	}

	if c.BlockRangeCap < 0 {
		return errors.New("JSON-RPC block range cap cannot be negative")
	}

	if c.HTTPTimeout < 0 {
		return errors.New("JSON-RPC HTTP timeout duration cannot be negative")
	}

	if c.HTTPIdleTimeout < 0 {
		return errors.New("JSON-RPC HTTP idle timeout duration cannot be negative")
	}

	// check for duplicates
	seenAPIs := make(map[string]bool)
	for _, api := range c.API {
		if seenAPIs[api] {
			return fmt.Errorf("repeated API namespace '%s'", api)
		}

		seenAPIs[api] = true
	}

	return nil
}

// DefaultTLSConfig returns the default TLS configuration
func DefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		CertificatePath: "",
		KeyPath:         "",
	}
}

// Validate returns an error if the TLS certificate and key file extensions are invalid.
func (c TLSConfig) Validate() error {
	certExt := path.Ext(c.CertificatePath)

	if c.CertificatePath != "" && certExt != ".pem" {
		return fmt.Errorf("invalid extension %s for certificate path %s, expected '.pem'", certExt, c.CertificatePath)
	}

	keyExt := path.Ext(c.KeyPath)

	if c.KeyPath != "" && keyExt != ".pem" {
		return fmt.Errorf("invalid extension %s for key path %s, expected '.pem'", keyExt, c.KeyPath)
	}

	return nil
}

// GetConfig returns a fully parsed Config object.
func GetConfig(v *viper.Viper) (Config, error) {
	cfg, err := config.GetConfig(v)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Config: cfg,
		EVM: EVMConfig{
			Tracer:         v.GetString("evm.tracer"),
			MaxTxGasWanted: v.GetUint64("evm.max-tx-gas-wanted"),
		},
		JSONRPC: JSONRPCConfig{
			Enable:             v.GetBool("json-rpc.enable"),
			API:                v.GetStringSlice("json-rpc.api"),
			Address:            v.GetString("json-rpc.address"),
			WsAddress:          v.GetString("json-rpc.ws-address"),
			GasCap:             v.GetUint64("json-rpc.gas-cap"),
			FilterCap:          v.GetInt32("json-rpc.filter-cap"),
			FeeHistoryCap:      v.GetInt32("json-rpc.feehistory-cap"),
			TxFeeCap:           v.GetFloat64("json-rpc.txfee-cap"),
			EVMTimeout:         v.GetDuration("json-rpc.evm-timeout"),
			LogsCap:            v.GetInt32("json-rpc.logs-cap"),
			BlockRangeCap:      v.GetInt32("json-rpc.block-range-cap"),
			HTTPTimeout:        v.GetDuration("json-rpc.http-timeout"),
			HTTPIdleTimeout:    v.GetDuration("json-rpc.http-idle-timeout"),
			MaxOpenConnections: v.GetInt("json-rpc.max-open-connections"),
			EnableIndexer:      v.GetBool("json-rpc.enable-indexer"),
			MetricsAddress:     v.GetString("json-rpc.metrics-address"),
		},
		TLS: TLSConfig{
			CertificatePath: v.GetString("tls.certificate-path"),
			KeyPath:         v.GetString("tls.key-path"),
		},
		Bundler: BundlerConfig{
			Enable:             viper.GetBool("bundler.enable"),
			Debug:              viper.GetBool("bundler.debug"),
			BeneficiaryAddress: viper.GetString("bundler.beneficiary"),
			EntryPointAddress:  viper.GetString("bundler.entrypoint"),
			MinBalance:         viper.GetString("bundler.min-balance"),
			SignerAddress:      viper.GetString("bundler.signer"),

			AutoBundle:            viper.GetBool("bundler.auto-bundle"),
			Whitelist:             viper.GetStringSlice("bundler.whitelist"),
			Blacklist:             viper.GetStringSlice("bundler.blacklist"),
			MaxBundleGas:          viper.GetUint64("bundler.max-bundle-gas"),
			MinStake:              viper.GetString("bundler.min-stake"),
			MinUnstakeDelay:       viper.GetUint64("bundler.min-unstake-delay"),
			AutoBundleInterval:    viper.GetUint64("bundler.auto-bundle-interval"),
			AutoBundleMempoolSize: viper.GetUint64("bundler.auto-bundle-mempool-size"),
		},
	}, nil
}

// ParseConfig retrieves the default environment configuration for the
// application.
func ParseConfig(v *viper.Viper) (*Config, error) {
	conf := DefaultConfig()
	err := v.Unmarshal(conf)

	return conf, err
}

// ValidateBasic returns an error any of the application configuration fields are invalid
func (c Config) ValidateBasic() error {
	if err := c.EVM.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid evm config value: %s", err.Error())
	}

	if err := c.JSONRPC.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid json-rpc config value: %s", err.Error())
	}

	if err := c.TLS.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid tls config value: %s", err.Error())
	}

	if err := c.Bundler.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid bundler config value: %s", err.Error())
	}

	return c.Config.ValidateBasic()
}
