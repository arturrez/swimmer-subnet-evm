package predeploy

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// StateDB is the interface for accessing EVM state
type StateDB interface {
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	SetCode(common.Address, []byte)

	SetNonce(common.Address, uint64)
	GetNonce(common.Address) uint64

	GetBalance(common.Address) *big.Int
	AddBalance(common.Address, *big.Int)
	SubBalance(common.Address, *big.Int)

	CreateAccount(common.Address)
	Exist(common.Address) bool
}

type PredeployContract struct {
	ContractAccessControlAddress        common.Address
	ContractConsensusAddress            common.Address
	ContractDefaultLogicAddress         common.Address
	AdminProxyStorageSlot               common.Hash
	LogicAddressStorageSlot             common.Hash
	GasLimitStorageSlot                 common.Hash
	GasPriceStorageSlot                 common.Hash
	RewardPoolStorageSlot               common.Hash
	DefaultLogicByteCode                []byte
	TransparentUpgradeableProxyByteCode []byte
}

func NewPredeployContract() *PredeployContract {
	return &PredeployContract{
		ContractAccessControlAddress,
		ContractConsensusAddress,
		ContractDefaultLogicAddress,
		AdminProxyStorageSlot,
		LogicAddressStorageSlot,
		GasLimitStorageSlot,
		GasPriceStorageSlot,
		RewardPoolStorageSlot,
		DefaultLogicByteCode,
		TransparentUpgradeableProxyByteCode,
	}
}

// Get gas price from consensus contract
func (p *PredeployContract) GetGasLimit(s StateDB) *big.Int {
	value := s.GetState(ContractConsensusAddress, GasLimitStorageSlot)
	return value.Big()
}

// Get gas price from consensus contract
func (p *PredeployContract) GetGasPrice(s StateDB) *big.Int {
	value := s.GetState(ContractConsensusAddress, GasPriceStorageSlot)
	return value.Big()
}

// Get reward pool address from consensus contract
func (p *PredeployContract) GetRewardPoolAddress(s StateDB) common.Address {
	value := s.GetState(ContractConsensusAddress, RewardPoolStorageSlot)
	return common.HexToAddress(value.Hex())
}

// Check if from is whitelisted user in the contract
func (p *PredeployContract) IsWhitelistedUser(s StateDB, scAddr common.Address, from common.Address) bool {
	storageSlot := common.BigToHash(big.NewInt(8)) // whitelisted users was stored at slot 8.
	scSlot := crypto.Keccak256Hash(append(common.LeftPadBytes(scAddr[:], 32)[:], storageSlot.Bytes()[:]...))
	fromSlot := crypto.Keccak256Hash(append(common.LeftPadBytes(from[:], 32)[:], scSlot.Bytes()[:]...))
	value := s.GetState(ContractAccessControlAddress, fromSlot)
	return value.Big().Cmp(big.NewInt(0)) == 1
}

// Check if func signature is accepted in cover fee mode
func (p *PredeployContract) IsWhitelistedFunc(s StateDB, scAddr common.Address, signature []byte) bool {
	storageSlot := common.BigToHash(big.NewInt(9)) // whitelisted funcs was stored at slot 9.
	scSlot := crypto.Keccak256Hash(append(common.LeftPadBytes(scAddr[:], 32)[:], storageSlot.Bytes()[:]...))
	signatureSlot := crypto.Keccak256Hash(append(common.RightPadBytes(signature[:], 32)[:], scSlot.Bytes()[:]...))
	value := s.GetState(ContractAccessControlAddress, signatureSlot)
	return value.Big().Cmp(big.NewInt(0)) == 1
}

func (p *PredeployContract) GetBlockedTimeOf(s StateDB, addr common.Address) uint64 {
	storageSlot := common.BigToHash(big.NewInt(2)) // mapping blocked time was stored at slot 2.
	addrKey := common.LeftPadBytes(addr[:], 32)
	addrSlot := crypto.Keccak256Hash(append(addrKey[:], storageSlot.Bytes()[:]...))
	value := s.GetState(ContractAccessControlAddress, addrSlot)
	return value.Big().Uint64()
}
