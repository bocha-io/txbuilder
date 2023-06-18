package txbuilder

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/bocha-io/garnet/x/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (t *TxBuilder) SendTransaction(address common.Address, privateKey *ecdsa.PrivateKey, message string, args ...interface{}) (common.Hash, error) {
	client, err := ethclient.Dial(t.endpoint)
	if err != nil {
		return common.Hash{}, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return common.Hash{}, err
	}

	value := big.NewInt(0)
	gasLimit := t.GetGasLimit(message)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}

	var data []byte
	if len(args) > 0 {
		data, err = t.worldABI.Pack(message, args...)
		if err != nil {
			return common.Hash{}, err
		}
	} else {
		data, err = t.worldABI.Pack(message)
		if err != nil {
			return common.Hash{}, err
		}
	}

	tx := types.NewTransaction(nonce, t.worldAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return common.Hash{}, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, err
	}

	logger.LogDebug(fmt.Sprintf("[backend] tx sent (%s) with hash: %s", message, signedTx.Hash().Hex()))

	return signedTx.Hash(), nil
}