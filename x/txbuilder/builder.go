package txbuilder

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	address common.Address
	ABI     abi.ABI
}

func NewContract(address string, abi abi.ABI) Contract {
	return Contract{
		address: common.HexToAddress(address),
		ABI:     abi,
	}
}

type TxBuilder struct {
	contracts map[string]Contract
	endpoint  string
	mnemonic  string

	customGasLimit  map[string]uint64
	defaultGasLimit uint64

	faucetPrivKey *ecdsa.PrivateKey

	currentNonce map[string]uint64
}

func NexTxBuilder(
	contracts map[string]Contract,
	endpoint string,
	mnemonic string,
	customGasLimit map[string]uint64,
	defaultGasLimit uint64,
	faucetPrivKey *ecdsa.PrivateKey,
) *TxBuilder {
	return &TxBuilder{
		contracts:       contracts,
		endpoint:        endpoint,
		mnemonic:        mnemonic,
		customGasLimit:  customGasLimit,
		defaultGasLimit: defaultGasLimit,
		faucetPrivKey:   faucetPrivKey,
		currentNonce:    map[string]uint64{},
	}
}

func (t *TxBuilder) InteractWithContract(
	contractName string,
	accountID int,
	message string,
	args ...interface{},
) (common.Hash, error) {
	wallet, account, err := GetWallet(t.mnemonic, accountID)
	if err != nil {
		return common.Hash{}, err
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return common.Hash{}, err
	}

	return t.SendTransaction(contractName, account.Address, privateKey, message, args...)
}

func (t *TxBuilder) FundAnAccount(accountID int) (common.Hash, error) {
	_, account, err := GetWallet(t.mnemonic, accountID)
	if err != nil {
		return common.Hash{}, err
	}
	// It sends 9 ETH
	return t.CallFaucet(account.Address.Hex(), big.NewInt(9000000000000000000))
}
