package txbuilder

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func (t *TxBuilder) GetGasLimit(method string) uint64 {
	v, ok := t.customGasLimit[method]
	if ok {
		return v
	}
	return t.defaultGasLimit
}

func (t *TxBuilder) TransactionReceipt(hash common.Hash) *types.Receipt {
	return t.rpcClient.TransactionReceipt(hash)
}

func (t *TxBuilder) WasTransactionSuccessful(hash common.Hash) (bool, error) {
	receipt := t.TransactionReceipt(hash)
	return receipt.Status == types.ReceiptStatusSuccessful, nil
}

func (t *TxBuilder) WasTxIncludedAndSuccessful(hash common.Hash) (bool, error) {
	retry := t.txCheckRetry
	for retry > 0 {
		res, err := t.WasTransactionSuccessful(hash)
		if err != nil {
			time.Sleep(t.txCheckWaitTime)
			retry--
			continue
		}
		return res, nil
	}

	return false, fmt.Errorf("checking tx timedout")
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
