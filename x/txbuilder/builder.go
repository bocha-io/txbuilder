package txbuilder

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type TxBuilder struct {
	worldAddress common.Address
	worldABI     abi.ABI

	endpoint string
	mnemonic string

	customGasLimit  map[string]uint64
	defaultGasLimit uint64

	faucetPrivKey *ecdsa.PrivateKey
}

func NexTxBuilder(
	worldAddress string,
	worldABI abi.ABI,
	endpoint string,
	mnemonic string,
	customGasLimit map[string]uint64,
	defaultGasLimit uint64,
	faucetPrivKey *ecdsa.PrivateKey,
) *TxBuilder {
	return &TxBuilder{
		worldAddress:    common.HexToAddress(worldAddress),
		worldABI:        worldABI,
		endpoint:        endpoint,
		mnemonic:        mnemonic,
		customGasLimit:  customGasLimit,
		defaultGasLimit: defaultGasLimit,
		faucetPrivKey:   faucetPrivKey,
	}
}

func (t *TxBuilder) InteractWithContract(
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

	return t.SendTransaction(account.Address, privateKey, message, args...)
}

func (t *TxBuilder) FoundAccount(accountID int) (common.Hash, error) {
	_, account, err := GetWallet(t.mnemonic, accountID)
	if err != nil {
		return common.Hash{}, err
	}
	// It sends 9 ETH
	return t.CallFaucet(account.Address.Hex(), big.NewInt(9000000000000000000))
}
