// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package relaying

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// RelayerManagerMetaData contains all meta data concerning the RelayerManager contract.
var RelayerManagerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"}],\"name\":\"ModuleCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"refundAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"refundToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signedHash\",\"type\":\"bytes32\"}],\"name\":\"TransactionExecuted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_module\",\"type\":\"address\"}],\"name\":\"addModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"}],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"recoverToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_methodId\",\"type\":\"bytes4\"}],\"name\":\"supportsStaticCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_isSupported\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"getRequiredSignatures\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"enumBaseModule.OwnerSignature\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_signatures\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_refundToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_refundAddress\",\"type\":\"address\"}],\"name\":\"execute\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"}],\"name\":\"getNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_signHash\",\"type\":\"bytes32\"}],\"name\":\"isExecutedTx\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"executed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"}],\"name\":\"getSession\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"key\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"expires\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// RelayerManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use RelayerManagerMetaData.ABI instead.
var RelayerManagerABI = RelayerManagerMetaData.ABI

// RelayerManager is an auto generated Go binding around an Ethereum contract.
type RelayerManager struct {
	RelayerManagerCaller     // Read-only binding to the contract
	RelayerManagerTransactor // Write-only binding to the contract
	RelayerManagerFilterer   // Log filterer for contract events
}

// RelayerManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type RelayerManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayerManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RelayerManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayerManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RelayerManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayerManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RelayerManagerSession struct {
	Contract     *RelayerManager   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RelayerManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RelayerManagerCallerSession struct {
	Contract *RelayerManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// RelayerManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RelayerManagerTransactorSession struct {
	Contract     *RelayerManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// RelayerManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type RelayerManagerRaw struct {
	Contract *RelayerManager // Generic contract binding to access the raw methods on
}

// RelayerManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RelayerManagerCallerRaw struct {
	Contract *RelayerManagerCaller // Generic read-only contract binding to access the raw methods on
}

// RelayerManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RelayerManagerTransactorRaw struct {
	Contract *RelayerManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRelayerManager creates a new instance of RelayerManager, bound to a specific deployed contract.
func NewRelayerManager(address common.Address, backend bind.ContractBackend) (*RelayerManager, error) {
	contract, err := bindRelayerManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RelayerManager{RelayerManagerCaller: RelayerManagerCaller{contract: contract}, RelayerManagerTransactor: RelayerManagerTransactor{contract: contract}, RelayerManagerFilterer: RelayerManagerFilterer{contract: contract}}, nil
}

// NewRelayerManagerCaller creates a new read-only instance of RelayerManager, bound to a specific deployed contract.
func NewRelayerManagerCaller(address common.Address, caller bind.ContractCaller) (*RelayerManagerCaller, error) {
	contract, err := bindRelayerManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RelayerManagerCaller{contract: contract}, nil
}

// NewRelayerManagerTransactor creates a new write-only instance of RelayerManager, bound to a specific deployed contract.
func NewRelayerManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*RelayerManagerTransactor, error) {
	contract, err := bindRelayerManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RelayerManagerTransactor{contract: contract}, nil
}

// NewRelayerManagerFilterer creates a new log filterer instance of RelayerManager, bound to a specific deployed contract.
func NewRelayerManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*RelayerManagerFilterer, error) {
	contract, err := bindRelayerManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RelayerManagerFilterer{contract: contract}, nil
}

// bindRelayerManager binds a generic wrapper to an already deployed contract.
func bindRelayerManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RelayerManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RelayerManager *RelayerManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RelayerManager.Contract.RelayerManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RelayerManager *RelayerManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RelayerManager.Contract.RelayerManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RelayerManager *RelayerManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RelayerManager.Contract.RelayerManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RelayerManager *RelayerManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RelayerManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RelayerManager *RelayerManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RelayerManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RelayerManager *RelayerManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RelayerManager.Contract.contract.Transact(opts, method, params...)
}

// GetNonce is a free data retrieval call binding the contract method 0x2d0335ab.
//
// Solidity: function getNonce(address _wallet) view returns(uint256 nonce)
func (_RelayerManager *RelayerManagerCaller) GetNonce(opts *bind.CallOpts, _wallet common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RelayerManager.contract.Call(opts, &out, "getNonce", _wallet)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNonce is a free data retrieval call binding the contract method 0x2d0335ab.
//
// Solidity: function getNonce(address _wallet) view returns(uint256 nonce)
func (_RelayerManager *RelayerManagerSession) GetNonce(_wallet common.Address) (*big.Int, error) {
	return _RelayerManager.Contract.GetNonce(&_RelayerManager.CallOpts, _wallet)
}

// GetNonce is a free data retrieval call binding the contract method 0x2d0335ab.
//
// Solidity: function getNonce(address _wallet) view returns(uint256 nonce)
func (_RelayerManager *RelayerManagerCallerSession) GetNonce(_wallet common.Address) (*big.Int, error) {
	return _RelayerManager.Contract.GetNonce(&_RelayerManager.CallOpts, _wallet)
}

// GetRequiredSignatures is a free data retrieval call binding the contract method 0x3b73d67f.
//
// Solidity: function getRequiredSignatures(address _wallet, bytes _data) view returns(uint256, uint8)
func (_RelayerManager *RelayerManagerCaller) GetRequiredSignatures(opts *bind.CallOpts, _wallet common.Address, _data []byte) (*big.Int, uint8, error) {
	var out []interface{}
	err := _RelayerManager.contract.Call(opts, &out, "getRequiredSignatures", _wallet, _data)

	if err != nil {
		return *new(*big.Int), *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return out0, out1, err

}

// GetRequiredSignatures is a free data retrieval call binding the contract method 0x3b73d67f.
//
// Solidity: function getRequiredSignatures(address _wallet, bytes _data) view returns(uint256, uint8)
func (_RelayerManager *RelayerManagerSession) GetRequiredSignatures(_wallet common.Address, _data []byte) (*big.Int, uint8, error) {
	return _RelayerManager.Contract.GetRequiredSignatures(&_RelayerManager.CallOpts, _wallet, _data)
}

// GetRequiredSignatures is a free data retrieval call binding the contract method 0x3b73d67f.
//
// Solidity: function getRequiredSignatures(address _wallet, bytes _data) view returns(uint256, uint8)
func (_RelayerManager *RelayerManagerCallerSession) GetRequiredSignatures(_wallet common.Address, _data []byte) (*big.Int, uint8, error) {
	return _RelayerManager.Contract.GetRequiredSignatures(&_RelayerManager.CallOpts, _wallet, _data)
}

// GetSession is a free data retrieval call binding the contract method 0x8c8e13b9.
//
// Solidity: function getSession(address _wallet) view returns(address key, uint64 expires)
func (_RelayerManager *RelayerManagerCaller) GetSession(opts *bind.CallOpts, _wallet common.Address) (struct {
	Key     common.Address
	Expires uint64
}, error) {
	var out []interface{}
	err := _RelayerManager.contract.Call(opts, &out, "getSession", _wallet)

	outstruct := new(struct {
		Key     common.Address
		Expires uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Key = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Expires = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// GetSession is a free data retrieval call binding the contract method 0x8c8e13b9.
//
// Solidity: function getSession(address _wallet) view returns(address key, uint64 expires)
func (_RelayerManager *RelayerManagerSession) GetSession(_wallet common.Address) (struct {
	Key     common.Address
	Expires uint64
}, error) {
	return _RelayerManager.Contract.GetSession(&_RelayerManager.CallOpts, _wallet)
}

// GetSession is a free data retrieval call binding the contract method 0x8c8e13b9.
//
// Solidity: function getSession(address _wallet) view returns(address key, uint64 expires)
func (_RelayerManager *RelayerManagerCallerSession) GetSession(_wallet common.Address) (struct {
	Key     common.Address
	Expires uint64
}, error) {
	return _RelayerManager.Contract.GetSession(&_RelayerManager.CallOpts, _wallet)
}

// IsExecutedTx is a free data retrieval call binding the contract method 0x60c0fdc0.
//
// Solidity: function isExecutedTx(address _wallet, bytes32 _signHash) view returns(bool executed)
func (_RelayerManager *RelayerManagerCaller) IsExecutedTx(opts *bind.CallOpts, _wallet common.Address, _signHash [32]byte) (bool, error) {
	var out []interface{}
	err := _RelayerManager.contract.Call(opts, &out, "isExecutedTx", _wallet, _signHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsExecutedTx is a free data retrieval call binding the contract method 0x60c0fdc0.
//
// Solidity: function isExecutedTx(address _wallet, bytes32 _signHash) view returns(bool executed)
func (_RelayerManager *RelayerManagerSession) IsExecutedTx(_wallet common.Address, _signHash [32]byte) (bool, error) {
	return _RelayerManager.Contract.IsExecutedTx(&_RelayerManager.CallOpts, _wallet, _signHash)
}

// IsExecutedTx is a free data retrieval call binding the contract method 0x60c0fdc0.
//
// Solidity: function isExecutedTx(address _wallet, bytes32 _signHash) view returns(bool executed)
func (_RelayerManager *RelayerManagerCallerSession) IsExecutedTx(_wallet common.Address, _signHash [32]byte) (bool, error) {
	return _RelayerManager.Contract.IsExecutedTx(&_RelayerManager.CallOpts, _wallet, _signHash)
}

// SupportsStaticCall is a free data retrieval call binding the contract method 0x25b50934.
//
// Solidity: function supportsStaticCall(bytes4 _methodId) view returns(bool _isSupported)
func (_RelayerManager *RelayerManagerCaller) SupportsStaticCall(opts *bind.CallOpts, _methodId [4]byte) (bool, error) {
	var out []interface{}
	err := _RelayerManager.contract.Call(opts, &out, "supportsStaticCall", _methodId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsStaticCall is a free data retrieval call binding the contract method 0x25b50934.
//
// Solidity: function supportsStaticCall(bytes4 _methodId) view returns(bool _isSupported)
func (_RelayerManager *RelayerManagerSession) SupportsStaticCall(_methodId [4]byte) (bool, error) {
	return _RelayerManager.Contract.SupportsStaticCall(&_RelayerManager.CallOpts, _methodId)
}

// SupportsStaticCall is a free data retrieval call binding the contract method 0x25b50934.
//
// Solidity: function supportsStaticCall(bytes4 _methodId) view returns(bool _isSupported)
func (_RelayerManager *RelayerManagerCallerSession) SupportsStaticCall(_methodId [4]byte) (bool, error) {
	return _RelayerManager.Contract.SupportsStaticCall(&_RelayerManager.CallOpts, _methodId)
}

// AddModule is a paid mutator transaction binding the contract method 0x5a1db8c4.
//
// Solidity: function addModule(address _wallet, address _module) returns()
func (_RelayerManager *RelayerManagerTransactor) AddModule(opts *bind.TransactOpts, _wallet common.Address, _module common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "addModule", _wallet, _module)
}

// AddModule is a paid mutator transaction binding the contract method 0x5a1db8c4.
//
// Solidity: function addModule(address _wallet, address _module) returns()
func (_RelayerManager *RelayerManagerSession) AddModule(_wallet common.Address, _module common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.AddModule(&_RelayerManager.TransactOpts, _wallet, _module)
}

// AddModule is a paid mutator transaction binding the contract method 0x5a1db8c4.
//
// Solidity: function addModule(address _wallet, address _module) returns()
func (_RelayerManager *RelayerManagerTransactorSession) AddModule(_wallet common.Address, _module common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.AddModule(&_RelayerManager.TransactOpts, _wallet, _module)
}

// Execute is a paid mutator transaction binding the contract method 0xe0724b6e.
//
// Solidity: function execute(address _wallet, bytes _data, uint256 _nonce, bytes _signatures, uint256 _gasPrice, uint256 _gasLimit, address _refundToken, address _refundAddress) returns(bool)
func (_RelayerManager *RelayerManagerTransactor) Execute(opts *bind.TransactOpts, _wallet common.Address, _data []byte, _nonce *big.Int, _signatures []byte, _gasPrice *big.Int, _gasLimit *big.Int, _refundToken common.Address, _refundAddress common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "execute", _wallet, _data, _nonce, _signatures, _gasPrice, _gasLimit, _refundToken, _refundAddress)
}

// Execute is a paid mutator transaction binding the contract method 0xe0724b6e.
//
// Solidity: function execute(address _wallet, bytes _data, uint256 _nonce, bytes _signatures, uint256 _gasPrice, uint256 _gasLimit, address _refundToken, address _refundAddress) returns(bool)
func (_RelayerManager *RelayerManagerSession) Execute(_wallet common.Address, _data []byte, _nonce *big.Int, _signatures []byte, _gasPrice *big.Int, _gasLimit *big.Int, _refundToken common.Address, _refundAddress common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.Execute(&_RelayerManager.TransactOpts, _wallet, _data, _nonce, _signatures, _gasPrice, _gasLimit, _refundToken, _refundAddress)
}

// Execute is a paid mutator transaction binding the contract method 0xe0724b6e.
//
// Solidity: function execute(address _wallet, bytes _data, uint256 _nonce, bytes _signatures, uint256 _gasPrice, uint256 _gasLimit, address _refundToken, address _refundAddress) returns(bool)
func (_RelayerManager *RelayerManagerTransactorSession) Execute(_wallet common.Address, _data []byte, _nonce *big.Int, _signatures []byte, _gasPrice *big.Int, _gasLimit *big.Int, _refundToken common.Address, _refundAddress common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.Execute(&_RelayerManager.TransactOpts, _wallet, _data, _nonce, _signatures, _gasPrice, _gasLimit, _refundToken, _refundAddress)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _wallet) returns()
func (_RelayerManager *RelayerManagerTransactor) Init(opts *bind.TransactOpts, _wallet common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "init", _wallet)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _wallet) returns()
func (_RelayerManager *RelayerManagerSession) Init(_wallet common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.Init(&_RelayerManager.TransactOpts, _wallet)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _wallet) returns()
func (_RelayerManager *RelayerManagerTransactorSession) Init(_wallet common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.Init(&_RelayerManager.TransactOpts, _wallet)
}

// RecoverToken is a paid mutator transaction binding the contract method 0x9be65a60.
//
// Solidity: function recoverToken(address _token) returns()
func (_RelayerManager *RelayerManagerTransactor) RecoverToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _RelayerManager.contract.Transact(opts, "recoverToken", _token)
}

// RecoverToken is a paid mutator transaction binding the contract method 0x9be65a60.
//
// Solidity: function recoverToken(address _token) returns()
func (_RelayerManager *RelayerManagerSession) RecoverToken(_token common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.RecoverToken(&_RelayerManager.TransactOpts, _token)
}

// RecoverToken is a paid mutator transaction binding the contract method 0x9be65a60.
//
// Solidity: function recoverToken(address _token) returns()
func (_RelayerManager *RelayerManagerTransactorSession) RecoverToken(_token common.Address) (*types.Transaction, error) {
	return _RelayerManager.Contract.RecoverToken(&_RelayerManager.TransactOpts, _token)
}

// RelayerManagerModuleCreatedIterator is returned from FilterModuleCreated and is used to iterate over the raw logs and unpacked data for ModuleCreated events raised by the RelayerManager contract.
type RelayerManagerModuleCreatedIterator struct {
	Event *RelayerManagerModuleCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayerManagerModuleCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerManagerModuleCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayerManagerModuleCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayerManagerModuleCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerManagerModuleCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerManagerModuleCreated represents a ModuleCreated event raised by the RelayerManager contract.
type RelayerManagerModuleCreated struct {
	Name [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterModuleCreated is a free log retrieval operation binding the contract event 0x3019c8fc80239e3dff8f781212ae2004839c2cb61d6c70acd279ac65392145df.
//
// Solidity: event ModuleCreated(bytes32 name)
func (_RelayerManager *RelayerManagerFilterer) FilterModuleCreated(opts *bind.FilterOpts) (*RelayerManagerModuleCreatedIterator, error) {

	logs, sub, err := _RelayerManager.contract.FilterLogs(opts, "ModuleCreated")
	if err != nil {
		return nil, err
	}
	return &RelayerManagerModuleCreatedIterator{contract: _RelayerManager.contract, event: "ModuleCreated", logs: logs, sub: sub}, nil
}

// WatchModuleCreated is a free log subscription operation binding the contract event 0x3019c8fc80239e3dff8f781212ae2004839c2cb61d6c70acd279ac65392145df.
//
// Solidity: event ModuleCreated(bytes32 name)
func (_RelayerManager *RelayerManagerFilterer) WatchModuleCreated(opts *bind.WatchOpts, sink chan<- *RelayerManagerModuleCreated) (event.Subscription, error) {

	logs, sub, err := _RelayerManager.contract.WatchLogs(opts, "ModuleCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerManagerModuleCreated)
				if err := _RelayerManager.contract.UnpackLog(event, "ModuleCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseModuleCreated is a log parse operation binding the contract event 0x3019c8fc80239e3dff8f781212ae2004839c2cb61d6c70acd279ac65392145df.
//
// Solidity: event ModuleCreated(bytes32 name)
func (_RelayerManager *RelayerManagerFilterer) ParseModuleCreated(log types.Log) (*RelayerManagerModuleCreated, error) {
	event := new(RelayerManagerModuleCreated)
	if err := _RelayerManager.contract.UnpackLog(event, "ModuleCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayerManagerRefundIterator is returned from FilterRefund and is used to iterate over the raw logs and unpacked data for Refund events raised by the RelayerManager contract.
type RelayerManagerRefundIterator struct {
	Event *RelayerManagerRefund // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayerManagerRefundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerManagerRefund)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayerManagerRefund)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayerManagerRefundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerManagerRefundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerManagerRefund represents a Refund event raised by the RelayerManager contract.
type RelayerManagerRefund struct {
	Wallet        common.Address
	RefundAddress common.Address
	RefundToken   common.Address
	RefundAmount  *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRefund is a free log retrieval operation binding the contract event 0x22edd2bbb0b0afbdcf90d91da8a5e2100f8d8f67cdc766dee1742e9a36d6add3.
//
// Solidity: event Refund(address indexed wallet, address indexed refundAddress, address refundToken, uint256 refundAmount)
func (_RelayerManager *RelayerManagerFilterer) FilterRefund(opts *bind.FilterOpts, wallet []common.Address, refundAddress []common.Address) (*RelayerManagerRefundIterator, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}
	var refundAddressRule []interface{}
	for _, refundAddressItem := range refundAddress {
		refundAddressRule = append(refundAddressRule, refundAddressItem)
	}

	logs, sub, err := _RelayerManager.contract.FilterLogs(opts, "Refund", walletRule, refundAddressRule)
	if err != nil {
		return nil, err
	}
	return &RelayerManagerRefundIterator{contract: _RelayerManager.contract, event: "Refund", logs: logs, sub: sub}, nil
}

// WatchRefund is a free log subscription operation binding the contract event 0x22edd2bbb0b0afbdcf90d91da8a5e2100f8d8f67cdc766dee1742e9a36d6add3.
//
// Solidity: event Refund(address indexed wallet, address indexed refundAddress, address refundToken, uint256 refundAmount)
func (_RelayerManager *RelayerManagerFilterer) WatchRefund(opts *bind.WatchOpts, sink chan<- *RelayerManagerRefund, wallet []common.Address, refundAddress []common.Address) (event.Subscription, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}
	var refundAddressRule []interface{}
	for _, refundAddressItem := range refundAddress {
		refundAddressRule = append(refundAddressRule, refundAddressItem)
	}

	logs, sub, err := _RelayerManager.contract.WatchLogs(opts, "Refund", walletRule, refundAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerManagerRefund)
				if err := _RelayerManager.contract.UnpackLog(event, "Refund", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRefund is a log parse operation binding the contract event 0x22edd2bbb0b0afbdcf90d91da8a5e2100f8d8f67cdc766dee1742e9a36d6add3.
//
// Solidity: event Refund(address indexed wallet, address indexed refundAddress, address refundToken, uint256 refundAmount)
func (_RelayerManager *RelayerManagerFilterer) ParseRefund(log types.Log) (*RelayerManagerRefund, error) {
	event := new(RelayerManagerRefund)
	if err := _RelayerManager.contract.UnpackLog(event, "Refund", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayerManagerTransactionExecutedIterator is returned from FilterTransactionExecuted and is used to iterate over the raw logs and unpacked data for TransactionExecuted events raised by the RelayerManager contract.
type RelayerManagerTransactionExecutedIterator struct {
	Event *RelayerManagerTransactionExecuted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayerManagerTransactionExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerManagerTransactionExecuted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayerManagerTransactionExecuted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayerManagerTransactionExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerManagerTransactionExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerManagerTransactionExecuted represents a TransactionExecuted event raised by the RelayerManager contract.
type RelayerManagerTransactionExecuted struct {
	Wallet     common.Address
	Success    bool
	ReturnData []byte
	SignedHash [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTransactionExecuted is a free log retrieval operation binding the contract event 0x7da4525a280527268ba2e963ee6c1b18f43c9507bcb1d2560f652ab17c76e90a.
//
// Solidity: event TransactionExecuted(address indexed wallet, bool indexed success, bytes returnData, bytes32 signedHash)
func (_RelayerManager *RelayerManagerFilterer) FilterTransactionExecuted(opts *bind.FilterOpts, wallet []common.Address, success []bool) (*RelayerManagerTransactionExecutedIterator, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _RelayerManager.contract.FilterLogs(opts, "TransactionExecuted", walletRule, successRule)
	if err != nil {
		return nil, err
	}
	return &RelayerManagerTransactionExecutedIterator{contract: _RelayerManager.contract, event: "TransactionExecuted", logs: logs, sub: sub}, nil
}

// WatchTransactionExecuted is a free log subscription operation binding the contract event 0x7da4525a280527268ba2e963ee6c1b18f43c9507bcb1d2560f652ab17c76e90a.
//
// Solidity: event TransactionExecuted(address indexed wallet, bool indexed success, bytes returnData, bytes32 signedHash)
func (_RelayerManager *RelayerManagerFilterer) WatchTransactionExecuted(opts *bind.WatchOpts, sink chan<- *RelayerManagerTransactionExecuted, wallet []common.Address, success []bool) (event.Subscription, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _RelayerManager.contract.WatchLogs(opts, "TransactionExecuted", walletRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerManagerTransactionExecuted)
				if err := _RelayerManager.contract.UnpackLog(event, "TransactionExecuted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransactionExecuted is a log parse operation binding the contract event 0x7da4525a280527268ba2e963ee6c1b18f43c9507bcb1d2560f652ab17c76e90a.
//
// Solidity: event TransactionExecuted(address indexed wallet, bool indexed success, bytes returnData, bytes32 signedHash)
func (_RelayerManager *RelayerManagerFilterer) ParseTransactionExecuted(log types.Log) (*RelayerManagerTransactionExecuted, error) {
	event := new(RelayerManagerTransactionExecuted)
	if err := _RelayerManager.contract.UnpackLog(event, "TransactionExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
