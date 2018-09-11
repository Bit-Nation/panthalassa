// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package documents

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NotaryABI is the input ABI used to generate the binding from.
const NotaryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"records\",\"outputs\":[{\"name\":\"notarisedData\",\"type\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_notarisedData\",\"type\":\"bytes\"}],\"name\":\"record\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_record\",\"type\":\"bytes\"}],\"name\":\"notarize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// NotaryBin is the compiled bytecode used for deploying new contracts.
const NotaryBin = `0x608060405234801561001057600080fd5b50610522806100206000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166301e64725811461005b578063e1112648146100f2578063fb1ace341461014b575b600080fd5b34801561006757600080fd5b506100736004356101a6565b6040518080602001838152602001828103825284818151815260200191508051906020019080838360005b838110156100b657818101518382015260200161009e565b50505050905090810190601f1680156100e35780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b3480156100fe57600080fd5b506040805160206004803580820135601f810184900484028501840190955284845261007394369492936024939284019190819084018382808284375094975061024a9650505050505050565b34801561015757600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526101a494369492936024939284019190819084018382808284375094975061037e9650505050505050565b005b6000602081815291815260409081902080548251601f60026000196101006001861615020190931692909204918201859004850281018501909352808352909283919083018282801561023a5780601f1061020f5761010080835404028352916020019161023a565b820191906000526020600020905b81548152906001019060200180831161021d57829003601f168201915b5050505050908060010154905082565b60606000610256610443565b600080856040518082805190602001908083835b602083106102895780518252601f19909201916020918201910161026a565b518151600019602094850361010090810a8201928316921993909316919091179092526040805196909401869003909520885287820198909852958101600020815181546060601f60026001841615909802909b019091169590950498890188900490970287018401825290860187815295969095879550935085929091508401828280156103595780601f1061032e57610100808354040283529160200191610359565b820191906000526020600020905b81548152906001019060200180831161033c57829003601f168201915b5050509183525050600191909101546020918201528151910151909590945092505050565b6000816040518082805190602001908083835b602083106103b05780518252601f199092019160209182019101610391565b51815160209384036101000a6000190180199092169116179052604080519290940182900390912060008181529182905292902060010154919450501591506103fa905057600080fd5b604080518082018252838152426020808301919091526000848152808252929092208151805192939192610431928492019061045b565b50602082015181600101559050505050565b60408051808201909152606081526000602082015290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061049c57805160ff19168380011785556104c9565b828001600101855582156104c9579182015b828111156104c95782518255916020019190600101906104ae565b506104d59291506104d9565b5090565b6104f391905b808211156104d557600081556001016104df565b905600a165627a7a72305820b7df842df5393139052e21c871c46429425f1394048c184e0e0e7b58843ea51b0029`

// DeployNotary deploys a new Ethereum contract, binding an instance of Notary to it.
func DeployNotary(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Notary, error) {
	parsed, err := abi.JSON(strings.NewReader(NotaryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NotaryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Notary{NotaryCaller: NotaryCaller{contract: contract}, NotaryTransactor: NotaryTransactor{contract: contract}, NotaryFilterer: NotaryFilterer{contract: contract}}, nil
}

// Notary is an auto generated Go binding around an Ethereum contract.
type Notary struct {
	NotaryCaller     // Read-only binding to the contract
	NotaryTransactor // Write-only binding to the contract
	NotaryFilterer   // Log filterer for contract events
}

// NotaryCaller is an auto generated read-only Go binding around an Ethereum contract.
type NotaryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NotaryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NotaryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NotaryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NotaryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NotarySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NotarySession struct {
	Contract     *Notary           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NotaryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NotaryCallerSession struct {
	Contract *NotaryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// NotaryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NotaryTransactorSession struct {
	Contract     *NotaryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NotaryRaw is an auto generated low-level Go binding around an Ethereum contract.
type NotaryRaw struct {
	Contract *Notary // Generic contract binding to access the raw methods on
}

// NotaryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NotaryCallerRaw struct {
	Contract *NotaryCaller // Generic read-only contract binding to access the raw methods on
}

// NotaryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NotaryTransactorRaw struct {
	Contract *NotaryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNotary creates a new instance of Notary, bound to a specific deployed contract.
func NewNotary(address common.Address, backend bind.ContractBackend) (*Notary, error) {
	contract, err := bindNotary(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Notary{NotaryCaller: NotaryCaller{contract: contract}, NotaryTransactor: NotaryTransactor{contract: contract}, NotaryFilterer: NotaryFilterer{contract: contract}}, nil
}

// NewNotaryCaller creates a new read-only instance of Notary, bound to a specific deployed contract.
func NewNotaryCaller(address common.Address, caller bind.ContractCaller) (*NotaryCaller, error) {
	contract, err := bindNotary(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NotaryCaller{contract: contract}, nil
}

// NewNotaryTransactor creates a new write-only instance of Notary, bound to a specific deployed contract.
func NewNotaryTransactor(address common.Address, transactor bind.ContractTransactor) (*NotaryTransactor, error) {
	contract, err := bindNotary(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NotaryTransactor{contract: contract}, nil
}

// NewNotaryFilterer creates a new log filterer instance of Notary, bound to a specific deployed contract.
func NewNotaryFilterer(address common.Address, filterer bind.ContractFilterer) (*NotaryFilterer, error) {
	contract, err := bindNotary(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NotaryFilterer{contract: contract}, nil
}

// bindNotary binds a generic wrapper to an already deployed contract.
func bindNotary(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NotaryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Notary *NotaryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Notary.Contract.NotaryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Notary *NotaryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Notary.Contract.NotaryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Notary *NotaryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Notary.Contract.NotaryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Notary *NotaryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Notary.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Notary *NotaryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Notary.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Notary *NotaryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Notary.Contract.contract.Transact(opts, method, params...)
}

// Record is a free data retrieval call binding the contract method 0xe1112648.
//
// Solidity: function record(_notarisedData bytes) constant returns(bytes, uint256)
func (_Notary *NotaryCaller) Record(opts *bind.CallOpts, _notarisedData []byte) ([]byte, *big.Int, error) {
	var (
		ret0 = new([]byte)
		ret1 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Notary.contract.Call(opts, out, "record", _notarisedData)
	return *ret0, *ret1, err
}

// Record is a free data retrieval call binding the contract method 0xe1112648.
//
// Solidity: function record(_notarisedData bytes) constant returns(bytes, uint256)
func (_Notary *NotarySession) Record(_notarisedData []byte) ([]byte, *big.Int, error) {
	return _Notary.Contract.Record(&_Notary.CallOpts, _notarisedData)
}

// Record is a free data retrieval call binding the contract method 0xe1112648.
//
// Solidity: function record(_notarisedData bytes) constant returns(bytes, uint256)
func (_Notary *NotaryCallerSession) Record(_notarisedData []byte) ([]byte, *big.Int, error) {
	return _Notary.Contract.Record(&_Notary.CallOpts, _notarisedData)
}

// Records is a free data retrieval call binding the contract method 0x01e64725.
//
// Solidity: function records( bytes32) constant returns(notarisedData bytes, timestamp uint256)
func (_Notary *NotaryCaller) Records(opts *bind.CallOpts, arg0 [32]byte) (struct {
	NotarisedData []byte
	Timestamp     *big.Int
}, error) {
	ret := new(struct {
		NotarisedData []byte
		Timestamp     *big.Int
	})
	out := ret
	err := _Notary.contract.Call(opts, out, "records", arg0)
	return *ret, err
}

// Records is a free data retrieval call binding the contract method 0x01e64725.
//
// Solidity: function records( bytes32) constant returns(notarisedData bytes, timestamp uint256)
func (_Notary *NotarySession) Records(arg0 [32]byte) (struct {
	NotarisedData []byte
	Timestamp     *big.Int
}, error) {
	return _Notary.Contract.Records(&_Notary.CallOpts, arg0)
}

// Records is a free data retrieval call binding the contract method 0x01e64725.
//
// Solidity: function records( bytes32) constant returns(notarisedData bytes, timestamp uint256)
func (_Notary *NotaryCallerSession) Records(arg0 [32]byte) (struct {
	NotarisedData []byte
	Timestamp     *big.Int
}, error) {
	return _Notary.Contract.Records(&_Notary.CallOpts, arg0)
}

// Notarize is a paid mutator transaction binding the contract method 0xfb1ace34.
//
// Solidity: function notarize(_record bytes) returns()
func (_Notary *NotaryTransactor) Notarize(opts *bind.TransactOpts, _record []byte) (*types.Transaction, error) {
	return _Notary.contract.Transact(opts, "notarize", _record)
}

// Notarize is a paid mutator transaction binding the contract method 0xfb1ace34.
//
// Solidity: function notarize(_record bytes) returns()
func (_Notary *NotarySession) Notarize(_record []byte) (*types.Transaction, error) {
	return _Notary.Contract.Notarize(&_Notary.TransactOpts, _record)
}

// Notarize is a paid mutator transaction binding the contract method 0xfb1ace34.
//
// Solidity: function notarize(_record bytes) returns()
func (_Notary *NotaryTransactorSession) Notarize(_record []byte) (*types.Transaction, error) {
	return _Notary.Contract.Notarize(&_Notary.TransactOpts, _record)
}

// NotaryMultiABI is the input ABI used to generate the binding from.
const NotaryMultiABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_notaryAddress\",\"type\":\"address\"},{\"name\":\"_firstRecord\",\"type\":\"bytes\"},{\"name\":\"_secondRecord\",\"type\":\"bytes\"}],\"name\":\"notarizeTwo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// NotaryMultiBin is the compiled bytecode used for deploying new contracts.
const NotaryMultiBin = `0x608060405234801561001057600080fd5b506102fb806100206000396000f3006080604052600436106100405763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416636ca484d78114610045575b600080fd5b34801561005157600080fd5b5060408051602060046024803582810135601f81018590048502860185019096528585526100f795833573ffffffffffffffffffffffffffffffffffffffff1695369560449491939091019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497506100f99650505050505050565b005b6040517ffb1ace34000000000000000000000000000000000000000000000000000000008152602060048201818152845160248401528451869373ffffffffffffffffffffffffffffffffffffffff85169363fb1ace34938893909283926044909101919085019080838360005b8381101561017f578181015183820152602001610167565b50505050905090810190601f1680156101ac5780820380516001836020036101000a031916815260200191505b5092505050600060405180830381600087803b1580156101cb57600080fd5b505af11580156101df573d6000803e3d6000fd5b50506040517ffb1ace3400000000000000000000000000000000000000000000000000000000815260206004820181815286516024840152865173ffffffffffffffffffffffffffffffffffffffff8716955063fb1ace349450879391928392604401919085019080838360005b8381101561026557818101518382015260200161024d565b50505050905090810190601f1680156102925780820380516001836020036101000a031916815260200191505b5092505050600060405180830381600087803b1580156102b157600080fd5b505af11580156102c5573d6000803e3d6000fd5b50505050505050505600a165627a7a723058209ccebce4442a22d9a140fb00f0dd8f77d860725161bd44fae7bc44697f557b250029`

// DeployNotaryMulti deploys a new Ethereum contract, binding an instance of NotaryMulti to it.
func DeployNotaryMulti(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NotaryMulti, error) {
	parsed, err := abi.JSON(strings.NewReader(NotaryMultiABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NotaryMultiBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NotaryMulti{NotaryMultiCaller: NotaryMultiCaller{contract: contract}, NotaryMultiTransactor: NotaryMultiTransactor{contract: contract}, NotaryMultiFilterer: NotaryMultiFilterer{contract: contract}}, nil
}

// NotaryMulti is an auto generated Go binding around an Ethereum contract.
type NotaryMulti struct {
	NotaryMultiCaller     // Read-only binding to the contract
	NotaryMultiTransactor // Write-only binding to the contract
	NotaryMultiFilterer   // Log filterer for contract events
}

// NotaryMultiCaller is an auto generated read-only Go binding around an Ethereum contract.
type NotaryMultiCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NotaryMultiTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NotaryMultiTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NotaryMultiFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NotaryMultiFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NotaryMultiSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NotaryMultiSession struct {
	Contract     *NotaryMulti      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NotaryMultiCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NotaryMultiCallerSession struct {
	Contract *NotaryMultiCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// NotaryMultiTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NotaryMultiTransactorSession struct {
	Contract     *NotaryMultiTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// NotaryMultiRaw is an auto generated low-level Go binding around an Ethereum contract.
type NotaryMultiRaw struct {
	Contract *NotaryMulti // Generic contract binding to access the raw methods on
}

// NotaryMultiCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NotaryMultiCallerRaw struct {
	Contract *NotaryMultiCaller // Generic read-only contract binding to access the raw methods on
}

// NotaryMultiTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NotaryMultiTransactorRaw struct {
	Contract *NotaryMultiTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNotaryMulti creates a new instance of NotaryMulti, bound to a specific deployed contract.
func NewNotaryMulti(address common.Address, backend bind.ContractBackend) (*NotaryMulti, error) {
	contract, err := bindNotaryMulti(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NotaryMulti{NotaryMultiCaller: NotaryMultiCaller{contract: contract}, NotaryMultiTransactor: NotaryMultiTransactor{contract: contract}, NotaryMultiFilterer: NotaryMultiFilterer{contract: contract}}, nil
}

// NewNotaryMultiCaller creates a new read-only instance of NotaryMulti, bound to a specific deployed contract.
func NewNotaryMultiCaller(address common.Address, caller bind.ContractCaller) (*NotaryMultiCaller, error) {
	contract, err := bindNotaryMulti(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NotaryMultiCaller{contract: contract}, nil
}

// NewNotaryMultiTransactor creates a new write-only instance of NotaryMulti, bound to a specific deployed contract.
func NewNotaryMultiTransactor(address common.Address, transactor bind.ContractTransactor) (*NotaryMultiTransactor, error) {
	contract, err := bindNotaryMulti(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NotaryMultiTransactor{contract: contract}, nil
}

// NewNotaryMultiFilterer creates a new log filterer instance of NotaryMulti, bound to a specific deployed contract.
func NewNotaryMultiFilterer(address common.Address, filterer bind.ContractFilterer) (*NotaryMultiFilterer, error) {
	contract, err := bindNotaryMulti(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NotaryMultiFilterer{contract: contract}, nil
}

// bindNotaryMulti binds a generic wrapper to an already deployed contract.
func bindNotaryMulti(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NotaryMultiABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NotaryMulti *NotaryMultiRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _NotaryMulti.Contract.NotaryMultiCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NotaryMulti *NotaryMultiRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NotaryMulti.Contract.NotaryMultiTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NotaryMulti *NotaryMultiRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NotaryMulti.Contract.NotaryMultiTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NotaryMulti *NotaryMultiCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _NotaryMulti.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NotaryMulti *NotaryMultiTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NotaryMulti.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NotaryMulti *NotaryMultiTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NotaryMulti.Contract.contract.Transact(opts, method, params...)
}

// NotarizeTwo is a paid mutator transaction binding the contract method 0x6ca484d7.
//
// Solidity: function notarizeTwo(_notaryAddress address, _firstRecord bytes, _secondRecord bytes) returns()
func (_NotaryMulti *NotaryMultiTransactor) NotarizeTwo(opts *bind.TransactOpts, _notaryAddress common.Address, _firstRecord []byte, _secondRecord []byte) (*types.Transaction, error) {
	return _NotaryMulti.contract.Transact(opts, "notarizeTwo", _notaryAddress, _firstRecord, _secondRecord)
}

// NotarizeTwo is a paid mutator transaction binding the contract method 0x6ca484d7.
//
// Solidity: function notarizeTwo(_notaryAddress address, _firstRecord bytes, _secondRecord bytes) returns()
func (_NotaryMulti *NotaryMultiSession) NotarizeTwo(_notaryAddress common.Address, _firstRecord []byte, _secondRecord []byte) (*types.Transaction, error) {
	return _NotaryMulti.Contract.NotarizeTwo(&_NotaryMulti.TransactOpts, _notaryAddress, _firstRecord, _secondRecord)
}

// NotarizeTwo is a paid mutator transaction binding the contract method 0x6ca484d7.
//
// Solidity: function notarizeTwo(_notaryAddress address, _firstRecord bytes, _secondRecord bytes) returns()
func (_NotaryMulti *NotaryMultiTransactorSession) NotarizeTwo(_notaryAddress common.Address, _firstRecord []byte, _secondRecord []byte) (*types.Transaction, error) {
	return _NotaryMulti.Contract.NotarizeTwo(&_NotaryMulti.TransactOpts, _notaryAddress, _firstRecord, _secondRecord)
}
