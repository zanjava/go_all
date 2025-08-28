// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ballot

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
	_ = abi.ConvertType
)

// BallotMetaData contains all meta data concerning the Ballot contract.
var BallotMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numProposals\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"enumBallot.Phase\",\"name\":\"x\",\"type\":\"uint8\"}],\"name\":\"changeState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWinner\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"state\",\"outputs\":[{\"internalType\":\"enumBallot.Phase\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"candidate\",\"type\":\"uint256\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523461017557604051601f6104c338819003918201601f19168301916001600160401b038311848410176101095780849260209460405283398101031261017557516003805460ff19168155811061011d575f80546001600160a01b0319163390811782558152600160205260408120600290555b81811061009b57600160ff196003541617600355604051610349908161017a8239f35b60405190602082016001600160401b03811183821017610109576040525f825260025491680100000000000000008310156101095760018301806002558310156100f55760019260025f5260205f20019051905501610078565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52604160045260245ffd5b60405162461bcd60e51b815260206004820152602a60248201527f63616e64696461746520636f756e742073686f756c64206e6f74206c657373206044820152697468616e20746872656560b01b6064820152608490fd5b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c80630121b93f146101fa578063268f1153146101a25780634420e4861461011d5780638e7ea5b2146100915763c19d93fb14610050575f80fd5b3461008d575f36600319011261008d5760ff600354166040516004821015610079576020918152f35b634e487b7160e01b5f52602160045260245ffd5b5f80fd5b3461008d575f36600319011261008d5760ff6003541660048110156100795760030361008d576002545f81815b8181106100ed575060011c116100d957602090604051908152f35b634e487b7160e01b5f52600160045260245ffd5b826100f782610291565b505411610107575b6001016100be565b9250905061011482610291565b505490826100ff565b3461008d57602036600319011261008d576004356001600160a01b0381169081900361008d5760ff6003541660048110156100795760010361008d5761016d60018060a01b035f541633146102bd565b805f52600160205260ff600160405f2001541661008d575f90815260016020819052604090912081815501805460ff19169055005b3461008d57602036600319011261008d57600435600481101561008d576101d360018060a01b035f541633146102bd565b6003549060ff8216600481101561007957811061008d5760ff169060ff1916176003555f80f35b3461008d57602036600319011261008d5760043560ff6003541660048110156100795760020361008d57335f52600160205260405f20600181019081549160ff83168015610285575b61008d5760016102619385600285015560ff19161790555491610291565b5080549182018092116102715755005b634e487b7160e01b5f52601160045260245ffd5b50600254841015610243565b6002548110156102a95760025f5260205f2001905f90565b634e487b7160e01b5f52603260045260245ffd5b156102c457565b60405162461bcd60e51b815260206004820152602160248201527f6f6e6c79206368616972506572736f6e2063616e206368616e676520737461746044820152606560f81b6064820152608490fdfea26469706673582212203bd3b4479e7c30ad784246ee3fdf04a63662d4154f375a33141f59567e15e28d64736f6c634300081e0033",
}

// BallotABI is the input ABI used to generate the binding from.
// Deprecated: Use BallotMetaData.ABI instead.
var BallotABI = BallotMetaData.ABI

// BallotBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BallotMetaData.Bin instead.
var BallotBin = BallotMetaData.Bin

// DeployBallot deploys a new Ethereum contract, binding an instance of Ballot to it.
func DeployBallot(auth *bind.TransactOpts, backend bind.ContractBackend, numProposals *big.Int) (common.Address, *types.Transaction, *Ballot, error) {
	parsed, err := BallotMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BallotBin), backend, numProposals)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ballot{BallotCaller: BallotCaller{contract: contract}, BallotTransactor: BallotTransactor{contract: contract}, BallotFilterer: BallotFilterer{contract: contract}}, nil
}

// Ballot is an auto generated Go binding around an Ethereum contract.
type Ballot struct {
	BallotCaller     // Read-only binding to the contract
	BallotTransactor // Write-only binding to the contract
	BallotFilterer   // Log filterer for contract events
}

// BallotCaller is an auto generated read-only Go binding around an Ethereum contract.
type BallotCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BallotTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BallotTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BallotFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BallotFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BallotSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BallotSession struct {
	Contract     *Ballot           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BallotCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BallotCallerSession struct {
	Contract *BallotCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BallotTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BallotTransactorSession struct {
	Contract     *BallotTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BallotRaw is an auto generated low-level Go binding around an Ethereum contract.
type BallotRaw struct {
	Contract *Ballot // Generic contract binding to access the raw methods on
}

// BallotCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BallotCallerRaw struct {
	Contract *BallotCaller // Generic read-only contract binding to access the raw methods on
}

// BallotTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BallotTransactorRaw struct {
	Contract *BallotTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBallot creates a new instance of Ballot, bound to a specific deployed contract.
func NewBallot(address common.Address, backend bind.ContractBackend) (*Ballot, error) {
	contract, err := bindBallot(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ballot{BallotCaller: BallotCaller{contract: contract}, BallotTransactor: BallotTransactor{contract: contract}, BallotFilterer: BallotFilterer{contract: contract}}, nil
}

// NewBallotCaller creates a new read-only instance of Ballot, bound to a specific deployed contract.
func NewBallotCaller(address common.Address, caller bind.ContractCaller) (*BallotCaller, error) {
	contract, err := bindBallot(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BallotCaller{contract: contract}, nil
}

// NewBallotTransactor creates a new write-only instance of Ballot, bound to a specific deployed contract.
func NewBallotTransactor(address common.Address, transactor bind.ContractTransactor) (*BallotTransactor, error) {
	contract, err := bindBallot(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BallotTransactor{contract: contract}, nil
}

// NewBallotFilterer creates a new log filterer instance of Ballot, bound to a specific deployed contract.
func NewBallotFilterer(address common.Address, filterer bind.ContractFilterer) (*BallotFilterer, error) {
	contract, err := bindBallot(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BallotFilterer{contract: contract}, nil
}

// bindBallot binds a generic wrapper to an already deployed contract.
func bindBallot(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BallotMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ballot *BallotRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ballot.Contract.BallotCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ballot *BallotRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ballot.Contract.BallotTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ballot *BallotRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ballot.Contract.BallotTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ballot *BallotCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ballot.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ballot *BallotTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ballot.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ballot *BallotTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ballot.Contract.contract.Transact(opts, method, params...)
}

// GetWinner is a free data retrieval call binding the contract method 0x8e7ea5b2.
//
// Solidity: function getWinner() view returns(uint256)
func (_Ballot *BallotCaller) GetWinner(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ballot.contract.Call(opts, &out, "getWinner")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetWinner is a free data retrieval call binding the contract method 0x8e7ea5b2.
//
// Solidity: function getWinner() view returns(uint256)
func (_Ballot *BallotSession) GetWinner() (*big.Int, error) {
	return _Ballot.Contract.GetWinner(&_Ballot.CallOpts)
}

// GetWinner is a free data retrieval call binding the contract method 0x8e7ea5b2.
//
// Solidity: function getWinner() view returns(uint256)
func (_Ballot *BallotCallerSession) GetWinner() (*big.Int, error) {
	return _Ballot.Contract.GetWinner(&_Ballot.CallOpts)
}

// State is a free data retrieval call binding the contract method 0xc19d93fb.
//
// Solidity: function state() view returns(uint8)
func (_Ballot *BallotCaller) State(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Ballot.contract.Call(opts, &out, "state")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// State is a free data retrieval call binding the contract method 0xc19d93fb.
//
// Solidity: function state() view returns(uint8)
func (_Ballot *BallotSession) State() (uint8, error) {
	return _Ballot.Contract.State(&_Ballot.CallOpts)
}

// State is a free data retrieval call binding the contract method 0xc19d93fb.
//
// Solidity: function state() view returns(uint8)
func (_Ballot *BallotCallerSession) State() (uint8, error) {
	return _Ballot.Contract.State(&_Ballot.CallOpts)
}

// ChangeState is a paid mutator transaction binding the contract method 0x268f1153.
//
// Solidity: function changeState(uint8 x) returns()
func (_Ballot *BallotTransactor) ChangeState(opts *bind.TransactOpts, x uint8) (*types.Transaction, error) {
	return _Ballot.contract.Transact(opts, "changeState", x)
}

// ChangeState is a paid mutator transaction binding the contract method 0x268f1153.
//
// Solidity: function changeState(uint8 x) returns()
func (_Ballot *BallotSession) ChangeState(x uint8) (*types.Transaction, error) {
	return _Ballot.Contract.ChangeState(&_Ballot.TransactOpts, x)
}

// ChangeState is a paid mutator transaction binding the contract method 0x268f1153.
//
// Solidity: function changeState(uint8 x) returns()
func (_Ballot *BallotTransactorSession) ChangeState(x uint8) (*types.Transaction, error) {
	return _Ballot.Contract.ChangeState(&_Ballot.TransactOpts, x)
}

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address voter) returns()
func (_Ballot *BallotTransactor) Register(opts *bind.TransactOpts, voter common.Address) (*types.Transaction, error) {
	return _Ballot.contract.Transact(opts, "register", voter)
}

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address voter) returns()
func (_Ballot *BallotSession) Register(voter common.Address) (*types.Transaction, error) {
	return _Ballot.Contract.Register(&_Ballot.TransactOpts, voter)
}

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address voter) returns()
func (_Ballot *BallotTransactorSession) Register(voter common.Address) (*types.Transaction, error) {
	return _Ballot.Contract.Register(&_Ballot.TransactOpts, voter)
}

// Vote is a paid mutator transaction binding the contract method 0x0121b93f.
//
// Solidity: function vote(uint256 candidate) returns()
func (_Ballot *BallotTransactor) Vote(opts *bind.TransactOpts, candidate *big.Int) (*types.Transaction, error) {
	return _Ballot.contract.Transact(opts, "vote", candidate)
}

// Vote is a paid mutator transaction binding the contract method 0x0121b93f.
//
// Solidity: function vote(uint256 candidate) returns()
func (_Ballot *BallotSession) Vote(candidate *big.Int) (*types.Transaction, error) {
	return _Ballot.Contract.Vote(&_Ballot.TransactOpts, candidate)
}

// Vote is a paid mutator transaction binding the contract method 0x0121b93f.
//
// Solidity: function vote(uint256 candidate) returns()
func (_Ballot *BallotTransactorSession) Vote(candidate *big.Int) (*types.Transaction, error) {
	return _Ballot.Contract.Vote(&_Ballot.TransactOpts, candidate)
}
