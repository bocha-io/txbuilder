package txbuilder

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/bocha-io/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (t *TxBuilder) SendTransaction(contractName string, address common.Address, privateKey *ecdsa.PrivateKey, value *big.Int, message string, args ...interface{}) (common.Hash, error) {
	var contractABI abi.ABI
	var contractAddress common.Address
	if v, ok := t.contracts[contractName]; ok {
		contractABI = v.ABI
		contractAddress = v.address
	} else {
		return common.Hash{}, fmt.Errorf("invalid contract name")
	}

	v, ok := t.currentNonce[address.Hex()]
	nonce := uint64(0)
	if ok {
		nonce = v
	} else {
		nonce = t.rpcClient.PendingNonceAt(address)
	}

	gasLimit := t.GetGasLimit(message)
	gasPrice := t.rpcClient.SuggestGasPrice()

	var data []byte
	var err error
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

	chainID := t.rpcClient.NetworkID()

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	t.rpcClient.SendTransaction(signedTx)

	t.currentNonce[address.Hex()] = nonce + 1

	logger.LogDebug(fmt.Sprintf("[backend] tx sent (%s) with hash: %s", message, signedTx.Hash().Hex()))

	return signedTx.Hash(), nil
}
