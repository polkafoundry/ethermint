// Copyright 2021 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/evmos/ethermint/blob/main/LICENSE
package config

// DefaultConfigTemplate defines the configuration template for the EVM RPC configuration
const DefaultConfigTemplate = `
###############################################################################
###                             EVM Configuration                           ###
###############################################################################

[evm]

# Tracer defines the 'vm.Tracer' type that the EVM will use when the node is run in
# debug mode. To enable tracing use the '--evm.tracer' flag when starting your node.
# Valid types are: json|struct|access_list|markdown
tracer = "{{ .EVM.Tracer }}"

# MaxTxGasWanted defines the gas wanted for each eth tx returned in ante handler in check tx mode.
max-tx-gas-wanted = {{ .EVM.MaxTxGasWanted }}

###############################################################################
###                           JSON RPC Configuration                        ###
###############################################################################

[json-rpc]

# Enable defines if the gRPC server should be enabled.
enable = {{ .JSONRPC.Enable }}

# Address defines the EVM RPC HTTP server address to bind to.
address = "{{ .JSONRPC.Address }}"

# Address defines the EVM WebSocket server address to bind to.
ws-address = "{{ .JSONRPC.WsAddress }}"

# API defines a list of JSON-RPC namespaces that should be enabled
# Example: "eth,txpool,personal,net,debug,web3"
api = "{{range $index, $elmt := .JSONRPC.API}}{{if $index}},{{$elmt}}{{else}}{{$elmt}}{{end}}{{end}}"

# GasCap sets a cap on gas that can be used in eth_call/estimateGas (0=infinite). Default: 25,000,000.
gas-cap = {{ .JSONRPC.GasCap }}

# EVMTimeout is the global timeout for eth_call. Default: 5s.
evm-timeout = "{{ .JSONRPC.EVMTimeout }}"

# TxFeeCap is the global tx-fee cap for send transaction. Default: 1eth.
txfee-cap = {{ .JSONRPC.TxFeeCap }}

# FilterCap sets the global cap for total number of filters that can be created
filter-cap = {{ .JSONRPC.FilterCap }}

# FeeHistoryCap sets the global cap for total number of blocks that can be fetched
feehistory-cap = {{ .JSONRPC.FeeHistoryCap }}

# LogsCap defines the max number of results can be returned from single 'eth_getLogs' query.
logs-cap = {{ .JSONRPC.LogsCap }}

# BlockRangeCap defines the max block range allowed for 'eth_getLogs' query.
block-range-cap = {{ .JSONRPC.BlockRangeCap }}

# HTTPTimeout is the read/write timeout of http json-rpc server.
http-timeout = "{{ .JSONRPC.HTTPTimeout }}"

# HTTPIdleTimeout is the idle timeout of http json-rpc server.
http-idle-timeout = "{{ .JSONRPC.HTTPIdleTimeout }}"

# AllowUnprotectedTxs restricts unprotected (non EIP155 signed) transactions to be submitted via
# the node's RPC when the global parameter is disabled.
allow-unprotected-txs = {{ .JSONRPC.AllowUnprotectedTxs }}

# MaxOpenConnections sets the maximum number of simultaneous connections
# for the server listener.
max-open-connections = {{ .JSONRPC.MaxOpenConnections }}

# EnableIndexer enables the custom transaction indexer for the EVM (ethereum transactions).
enable-indexer = {{ .JSONRPC.EnableIndexer }}

# MetricsAddress defines the EVM Metrics server address to bind to. Pass --metrics in CLI to enable
# Prometheus metrics path: /debug/metrics/prometheus
metrics-address = "{{ .JSONRPC.MetricsAddress }}"

###############################################################################
###                             TLS Configuration                           ###
###############################################################################

[tls]

# Certificate path defines the cert.pem file path for the TLS configuration.
certificate-path = "{{ .TLS.CertificatePath }}"

# Key path defines the key.pem file path for the TLS configuration.
key-path = "{{ .TLS.KeyPath }}"

###############################################################################
###                          Relayer Configuration                          ###
###############################################################################

[relayer]

# Enable defines if node should be run as a relayer.
enable = {{ .Relayer.Enable }}

# Address defines address of RelayerManager contract
address = "{{ .Relayer.Address }}"

# SenderAddresses defines addresses of keys which are used to sign relay transactions.
sender-addresses = "{{range $index, $elmt := .Relayer.SenderAddresses}}{{if $index}},{{$elmt}}{{else}}{{$elmt}}{{end}}{{end}}"

# RefundAddresses defines addresses that receive refund tokens from relay transactions.
refund-addresses = "{{range $index, $elmt := .Relayer.RefundAddresses}}{{if $index}},{{$elmt}}{{else}}{{$elmt}}{{end}}{{end}}"

# RefundTokens defines addresses of tokens that can be used when performing refund in relay transactions.
# Zero address token means the native token of the chain.
refund-tokens = "{{range $index, $elmt := .Relayer.RefundTokens}}{{if $index}},{{$elmt}}{{else}}{{$elmt}}{{end}}{{end}}"

# MinGasPrices defines the minimum gas prices corresponding with the refund tokens that a relayer is willing to accept
# for processing a relay call.
min-gas-prices = "{{range $index, $elmt := .Relayer.MinGasPrices}}{{if $index}},{{$elmt}}{{else}}{{$elmt}}{{end}}{{end}}"

# GasMultiplier is the multiplier when sending relay transaction after estimating gas.
# gas = gasMultiplier * estimatedGas
gas-multiplier = "{{ .Relayer.GasMultiplier }}"
`
