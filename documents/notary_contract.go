// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package documents

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// NotaryABI is the input ABI used to generate the binding from.
const NotaryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"records\",\"outputs\":[{\"name\":\"notarisedData\",\"type\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setNotarisationFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"notarisationFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_notarisedData\",\"type\":\"bytes\"}],\"name\":\"record\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_record\",\"type\":\"bytes\"}],\"name\":\"notarize\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// NotaryBin is the compiled bytecode used for deploying new contracts.
const NotaryBin = `0x608060405234801561001057600080fd5b506040516020806107e5833981016040525160008054600160a060020a03909216600160a060020a0319928316331790921691909117905561078e806100576000396000f30060806040526004361061008d5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166301e647258114610092578063715018a6146101295780638da5cb5b14610140578063c0496e5714610171578063c9d3a88514610189578063e1112648146101b0578063f2fde38b14610209578063fb1ace341461022a575b600080fd5b34801561009e57600080fd5b506100aa600435610276565b6040518080602001838152602001828103825284818151815260200191508051906020019080838360005b838110156100ed5781810151838201526020016100d5565b50505050905090810190601f16801561011a5780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b34801561013557600080fd5b5061013e610318565b005b34801561014c57600080fd5b50610155610384565b60408051600160a060020a039092168252519081900360200190f35b34801561017d57600080fd5b5061013e600435610393565b34801561019557600080fd5b5061019e6103af565b60408051918252519081900360200190f35b3480156101bc57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100aa9436949293602493928401919081908401838280828437509497506103b59650505050505050565b34801561021557600080fd5b5061013e600160a060020a03600435166104ea565b6040805160206004803580820135601f810184900484028501840190955284845261013e94369492936024939284019190819084018382808284375094975061050d9650505050505050565b60016020818152600092835260409283902080548451600294821615610100026000190190911693909304601f81018390048302840183019094528383529283918301828280156103085780601f106102dd57610100808354040283529160200191610308565b820191906000526020600020905b8154815290600101906020018083116102eb57829003601f168201915b5050505050908060010154905082565b600054600160a060020a0316331461032f57600080fd5b60008054604051600160a060020a03909116917ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482091a26000805473ffffffffffffffffffffffffffffffffffffffff19169055565b600054600160a060020a031681565b600054600160a060020a031633146103aa57600080fd5b600255565b60025481565b606060006103c16106af565b60016000856040518082805190602001908083835b602083106103f55780518252601f1990920191602091820191016103d6565b518151600019602094850361010090810a8201928316921993909316919091179092526040805196909401869003909520885287820198909852958101600020815181546060601f60026001841615909802909b019091169590950498890188900490970287018401825290860187815295969095879550935085929091508401828280156104c55780601f1061049a576101008083540402835291602001916104c5565b820191906000526020600020905b8154815290600101906020018083116104a857829003601f168201915b5050509183525050600191909101546020918201528151910151909590945092505050565b600054600160a060020a0316331461050157600080fd5b61050a81610632565b50565b60025460009034101561051f57600080fd5b816040518082805190602001908083835b6020831061054f5780518252601f199092019160209182019101610530565b51815160209384036101000a6000190180199092169116179052604080519290940182900390912060008181526001928390529390932001549194505015915061059a905057600080fd5b600054600160a060020a0316156105e85760008054604051600160a060020a0390911691303180156108fc02929091818181858888f193505050501580156105e6573d6000803e3d6000fd5b505b6040805180820182528381524260208083019190915260008481526001825292909220815180519293919261062092849201906106c7565b50602082015181600101559050505050565b600160a060020a038116151561064757600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b60408051808201909152606081526000602082015290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061070857805160ff1916838001178555610735565b82800160010185558215610735579182015b8281111561073557825182559160200191906001019061071a565b50610741929150610745565b5090565b61075f91905b80821115610741576000815560010161074b565b905600a165627a7a7230582001b0252bf7f25e941fd126beef7b2b40c7c9922b862a8bd543807bd32f2cf4900029`

// DeployNotary deploys a new Ethereum contract, binding an instance of Notary to it.
func DeployNotary(auth *bind.TransactOpts, backend bind.ContractBackend, _owner common.Address) (common.Address, *types.Transaction, *Notary, error) {
	parsed, err := abi.JSON(strings.NewReader(NotaryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NotaryBin), backend, _owner)
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

// NotarisationFee is a free data retrieval call binding the contract method 0xc9d3a885.
//
// Solidity: function notarisationFee() constant returns(uint256)
func (_Notary *NotaryCaller) NotarisationFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Notary.contract.Call(opts, out, "notarisationFee")
	return *ret0, err
}

// NotarisationFee is a free data retrieval call binding the contract method 0xc9d3a885.
//
// Solidity: function notarisationFee() constant returns(uint256)
func (_Notary *NotarySession) NotarisationFee() (*big.Int, error) {
	return _Notary.Contract.NotarisationFee(&_Notary.CallOpts)
}

// NotarisationFee is a free data retrieval call binding the contract method 0xc9d3a885.
//
// Solidity: function notarisationFee() constant returns(uint256)
func (_Notary *NotaryCallerSession) NotarisationFee() (*big.Int, error) {
	return _Notary.Contract.NotarisationFee(&_Notary.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Notary *NotaryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Notary.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Notary *NotarySession) Owner() (common.Address, error) {
	return _Notary.Contract.Owner(&_Notary.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Notary *NotaryCallerSession) Owner() (common.Address, error) {
	return _Notary.Contract.Owner(&_Notary.CallOpts)
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

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Notary *NotaryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Notary.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Notary *NotarySession) RenounceOwnership() (*types.Transaction, error) {
	return _Notary.Contract.RenounceOwnership(&_Notary.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Notary *NotaryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Notary.Contract.RenounceOwnership(&_Notary.TransactOpts)
}

// SetNotarisationFee is a paid mutator transaction binding the contract method 0xc0496e57.
//
// Solidity: function setNotarisationFee(_fee uint256) returns()
func (_Notary *NotaryTransactor) SetNotarisationFee(opts *bind.TransactOpts, _fee *big.Int) (*types.Transaction, error) {
	return _Notary.contract.Transact(opts, "setNotarisationFee", _fee)
}

// SetNotarisationFee is a paid mutator transaction binding the contract method 0xc0496e57.
//
// Solidity: function setNotarisationFee(_fee uint256) returns()
func (_Notary *NotarySession) SetNotarisationFee(_fee *big.Int) (*types.Transaction, error) {
	return _Notary.Contract.SetNotarisationFee(&_Notary.TransactOpts, _fee)
}

// SetNotarisationFee is a paid mutator transaction binding the contract method 0xc0496e57.
//
// Solidity: function setNotarisationFee(_fee uint256) returns()
func (_Notary *NotaryTransactorSession) SetNotarisationFee(_fee *big.Int) (*types.Transaction, error) {
	return _Notary.Contract.SetNotarisationFee(&_Notary.TransactOpts, _fee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Notary *NotaryTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Notary.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Notary *NotarySession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Notary.Contract.TransferOwnership(&_Notary.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Notary *NotaryTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Notary.Contract.TransferOwnership(&_Notary.TransactOpts, _newOwner)
}

// NotaryOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the Notary contract.
type NotaryOwnershipRenouncedIterator struct {
	Event *NotaryOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *NotaryOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NotaryOwnershipRenounced)
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
		it.Event = new(NotaryOwnershipRenounced)
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
func (it *NotaryOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NotaryOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NotaryOwnershipRenounced represents a OwnershipRenounced event raised by the Notary contract.
type NotaryOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Notary *NotaryFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*NotaryOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Notary.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NotaryOwnershipRenouncedIterator{contract: _Notary.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Notary *NotaryFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *NotaryOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Notary.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NotaryOwnershipRenounced)
				if err := _Notary.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// NotaryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Notary contract.
type NotaryOwnershipTransferredIterator struct {
	Event *NotaryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *NotaryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NotaryOwnershipTransferred)
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
		it.Event = new(NotaryOwnershipTransferred)
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
func (it *NotaryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NotaryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NotaryOwnershipTransferred represents a OwnershipTransferred event raised by the Notary contract.
type NotaryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Notary *NotaryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NotaryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Notary.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NotaryOwnershipTransferredIterator{contract: _Notary.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Notary *NotaryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NotaryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Notary.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NotaryOwnershipTransferred)
				if err := _Notary.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// NotaryMultiABI is the input ABI used to generate the binding from.
const NotaryMultiABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_firstRecord\",\"type\":\"bytes\"},{\"name\":\"_secondRecord\",\"type\":\"bytes\"}],\"name\":\"notarizeTwo\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"notaryFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"notary\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_notary\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// NotaryMultiBin is the compiled bytecode used for deploying new contracts.
const NotaryMultiBin = `0x608060405234801561001057600080fd5b5060405160208061047d833981016040525160008054600160a060020a03909216600160a060020a031990921691909117905561042b806100526000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166309267785811461005b578063835c853b146100e75780639d54c79d1461010e575b600080fd5b6040805160206004803580820135601f81018490048402850184019095528484526100e594369492936024939284019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a99988101979196509182019450925082915084018382808284375094975061014c9650505050505050565b005b3480156100f357600080fd5b506100fc610328565b60408051918252519081900360200190f35b34801561011a57600080fd5b506101236103e3565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b600080546040517ffb1ace3400000000000000000000000000000000000000000000000000000000815260206004820181815286516024840152865173ffffffffffffffffffffffffffffffffffffffff9094169463fb1ace34948894929384936044019290860191908190849084905b838110156101d55781810151838201526020016101bd565b50505050905090810190601f1680156102025780820380516001836020036101000a031916815260200191505b5092505050600060405180830381600087803b15801561022157600080fd5b505af1158015610235573d6000803e3d6000fd5b5050600080546040517ffb1ace3400000000000000000000000000000000000000000000000000000000815260206004820181815287516024840152875173ffffffffffffffffffffffffffffffffffffffff909416965063fb1ace349550879490938493604401928601918190849084905b838110156102c05781810151838201526020016102a8565b50505050905090810190601f1680156102ed5780820380516001836020036101000a031916815260200191505b5092505050600060405180830381600087803b15801561030c57600080fd5b505af1158015610320573d6000803e3d6000fd5b505050505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c9d3a8856040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b1580156103af57600080fd5b505af11580156103c3573d6000803e3d6000fd5b505050506040513d60208110156103d957600080fd5b5051600202905090565b60005473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a7230582069c7b45b3c6de0d2511a2b2d5d56868ed00f9c9ba5582d5907c57c67f90b4f620029`

// DeployNotaryMulti deploys a new Ethereum contract, binding an instance of NotaryMulti to it.
func DeployNotaryMulti(auth *bind.TransactOpts, backend bind.ContractBackend, _notary common.Address) (common.Address, *types.Transaction, *NotaryMulti, error) {
	parsed, err := abi.JSON(strings.NewReader(NotaryMultiABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NotaryMultiBin), backend, _notary)
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

// Notary is a free data retrieval call binding the contract method 0x9d54c79d.
//
// Solidity: function notary() constant returns(address)
func (_NotaryMulti *NotaryMultiCaller) Notary(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _NotaryMulti.contract.Call(opts, out, "notary")
	return *ret0, err
}

// Notary is a free data retrieval call binding the contract method 0x9d54c79d.
//
// Solidity: function notary() constant returns(address)
func (_NotaryMulti *NotaryMultiSession) Notary() (common.Address, error) {
	return _NotaryMulti.Contract.Notary(&_NotaryMulti.CallOpts)
}

// Notary is a free data retrieval call binding the contract method 0x9d54c79d.
//
// Solidity: function notary() constant returns(address)
func (_NotaryMulti *NotaryMultiCallerSession) Notary() (common.Address, error) {
	return _NotaryMulti.Contract.Notary(&_NotaryMulti.CallOpts)
}

// NotaryFee is a free data retrieval call binding the contract method 0x835c853b.
//
// Solidity: function notaryFee() constant returns(uint256)
func (_NotaryMulti *NotaryMultiCaller) NotaryFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _NotaryMulti.contract.Call(opts, out, "notaryFee")
	return *ret0, err
}

// NotaryFee is a free data retrieval call binding the contract method 0x835c853b.
//
// Solidity: function notaryFee() constant returns(uint256)
func (_NotaryMulti *NotaryMultiSession) NotaryFee() (*big.Int, error) {
	return _NotaryMulti.Contract.NotaryFee(&_NotaryMulti.CallOpts)
}

// NotaryFee is a free data retrieval call binding the contract method 0x835c853b.
//
// Solidity: function notaryFee() constant returns(uint256)
func (_NotaryMulti *NotaryMultiCallerSession) NotaryFee() (*big.Int, error) {
	return _NotaryMulti.Contract.NotaryFee(&_NotaryMulti.CallOpts)
}

// NotarizeTwo is a paid mutator transaction binding the contract method 0x09267785.
//
// Solidity: function notarizeTwo(_firstRecord bytes, _secondRecord bytes) returns()
func (_NotaryMulti *NotaryMultiTransactor) NotarizeTwo(opts *bind.TransactOpts, _firstRecord []byte, _secondRecord []byte) (*types.Transaction, error) {
	return _NotaryMulti.contract.Transact(opts, "notarizeTwo", _firstRecord, _secondRecord)
}

// NotarizeTwo is a paid mutator transaction binding the contract method 0x09267785.
//
// Solidity: function notarizeTwo(_firstRecord bytes, _secondRecord bytes) returns()
func (_NotaryMulti *NotaryMultiSession) NotarizeTwo(_firstRecord []byte, _secondRecord []byte) (*types.Transaction, error) {
	return _NotaryMulti.Contract.NotarizeTwo(&_NotaryMulti.TransactOpts, _firstRecord, _secondRecord)
}

// NotarizeTwo is a paid mutator transaction binding the contract method 0x09267785.
//
// Solidity: function notarizeTwo(_firstRecord bytes, _secondRecord bytes) returns()
func (_NotaryMulti *NotaryMultiTransactorSession) NotarizeTwo(_firstRecord []byte, _secondRecord []byte) (*types.Transaction, error) {
	return _NotaryMulti.Contract.NotarizeTwo(&_NotaryMulti.TransactOpts, _firstRecord, _secondRecord)
}

// OwnableABI is the input ABI used to generate the binding from.
const OwnableABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// OwnableBin is the compiled bytecode used for deploying new contracts.
const OwnableBin = `0x608060405234801561001057600080fd5b5060008054600160a060020a0319163317905561020b806100326000396000f3006080604052600436106100565763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663715018a6811461005b5780638da5cb5b14610072578063f2fde38b146100a3575b600080fd5b34801561006757600080fd5b506100706100c4565b005b34801561007e57600080fd5b50610087610130565b60408051600160a060020a039092168252519081900360200190f35b3480156100af57600080fd5b50610070600160a060020a036004351661013f565b600054600160a060020a031633146100db57600080fd5b60008054604051600160a060020a03909116917ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482091a26000805473ffffffffffffffffffffffffffffffffffffffff19169055565b600054600160a060020a031681565b600054600160a060020a0316331461015657600080fd5b61015f81610162565b50565b600160a060020a038116151561017757600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a72305820939090677929d52a8b7fcd09005569e1345b68565d3dee73b40c555b2f7d4d300029`

// DeployOwnable deploys a new Ethereum contract, binding an instance of Ownable to it.
func DeployOwnable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Ownable, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// Ownable is an auto generated Go binding around an Ethereum contract.
type Ownable struct {
	OwnableCaller     // Read-only binding to the contract
	OwnableTransactor // Write-only binding to the contract
	OwnableFilterer   // Log filterer for contract events
}

// OwnableCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnableSession struct {
	Contract     *Ownable          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnableCallerSession struct {
	Contract *OwnableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// OwnableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnableTransactorSession struct {
	Contract     *OwnableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OwnableRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnableRaw struct {
	Contract *Ownable // Generic contract binding to access the raw methods on
}

// OwnableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnableCallerRaw struct {
	Contract *OwnableCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnableTransactorRaw struct {
	Contract *OwnableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnable creates a new instance of Ownable, bound to a specific deployed contract.
func NewOwnable(address common.Address, backend bind.ContractBackend) (*Ownable, error) {
	contract, err := bindOwnable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// NewOwnableCaller creates a new read-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableCaller(address common.Address, caller bind.ContractCaller) (*OwnableCaller, error) {
	contract, err := bindOwnable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableCaller{contract: contract}, nil
}

// NewOwnableTransactor creates a new write-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableTransactor, error) {
	contract, err := bindOwnable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableTransactor{contract: contract}, nil
}

// NewOwnableFilterer creates a new log filterer instance of Ownable, bound to a specific deployed contract.
func NewOwnableFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableFilterer, error) {
	contract, err := bindOwnable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableFilterer{contract: contract}, nil
}

// bindOwnable binds a generic wrapper to an already deployed contract.
func bindOwnable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.OwnableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Ownable.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableCallerSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Ownable *OwnableTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Ownable *OwnableSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Ownable *OwnableTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, _newOwner)
}

// OwnableOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the Ownable contract.
type OwnableOwnershipRenouncedIterator struct {
	Event *OwnableOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *OwnableOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableOwnershipRenounced)
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
		it.Event = new(OwnableOwnershipRenounced)
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
func (it *OwnableOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableOwnershipRenounced represents a OwnershipRenounced event raised by the Ownable contract.
type OwnableOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Ownable *OwnableFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*OwnableOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Ownable.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableOwnershipRenouncedIterator{contract: _Ownable.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Ownable *OwnableFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *OwnableOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Ownable.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableOwnershipRenounced)
				if err := _Ownable.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// OwnableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ownable contract.
type OwnableOwnershipTransferredIterator struct {
	Event *OwnableOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OwnableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableOwnershipTransferred)
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
		it.Event = new(OwnableOwnershipTransferred)
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
func (it *OwnableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableOwnershipTransferred represents a OwnershipTransferred event raised by the Ownable contract.
type OwnableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Ownable *OwnableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableOwnershipTransferredIterator{contract: _Ownable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Ownable *OwnableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableOwnershipTransferred)
				if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
