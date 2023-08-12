package txbuilder

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/bocha-io/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (t *TxBuilder) SendTransaction(contractName string, address common.Address, privateKey *ecdsa.PrivateKey, message string, args ...interface{}) (common.Hash, error) {
	client, err := ethclient.Dial(t.endpoint)
	if err != nil {
		return common.Hash{}, err
	}

	var contractABI abi.ABI
	var contractAddress common.Address
	if v, ok := t.contracts[contractName]; ok {
		contractABI = v.ABI
		contractAddress = v.address

	} else {
		return common.Hash{}, fmt.Errorf("invalid contract name")
	}

	v, ok := t.currentNonce[address.Hash().Hex()]
	nonce := uint64(0)
	if ok {
		nonce = v
	} else {
		nonce, err = client.PendingNonceAt(context.Background(), address)
		if err != nil {
			return common.Hash{}, err
		}
	}

	value := big.NewInt(0)
	gasLimit := t.GetGasLimit(message)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}

	var data []byte
	if len(args) > 0 {
		data, err = contractABI.Pack(message, args...)
		if err != nil {
			return common.Hash{}, err
		}
	} else {
		data, err = contractABI.Pack(message)
		if err != nil {
			return common.Hash{}, err
		}
	}

	tx := types.NewTransaction(nonce, contractAddress, value, gasLimit, gasPrice, data)
	logger.LogDebug(fmt.Sprintf("[backend] creating tx (%s) with nonce: %d", message, nonce))

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

	t.currentNonce[address.Hash().Hex()] = nonce + 1

	logger.LogDebug(fmt.Sprintf("[backend] tx sent (%s) with hash: %s", message, signedTx.Hash().Hex()))

	return signedTx.Hash(), nil
}
