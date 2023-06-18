package txbuilder

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func (t *TxBuilder) GetGasLimit(method string) uint64 {
	v, ok := t.customGasLimit[method]
	if ok {
		return v
	}
	return t.defaultGasLimit
}

func (t *TxBuilder) TransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	client, err := ethclient.Dial(t.endpoint)
	if err != nil {
		return nil, err
	}
	return client.TransactionReceipt(context.Background(), hash)
}

func (t *TxBuilder) WasTransactionSuccessful(hash common.Hash) (bool, error) {
	receipt, err := t.TransactionReceipt(hash)
	if err != nil {
		return false, err
	}
	return receipt.Status == types.ReceiptStatusSuccessful, nil
}

func GetWallet(mnemonic string, accountID int) (*hdwallet.Wallet, accounts.Account, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, accounts.Account{}, err
	}

	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", accountID))
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, accounts.Account{}, err
	}
	return wallet, account, nil
}
