package txbuilder

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/bocha-io/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var PrivateKeysAnvil = [3]*ecdsa.PrivateKey{}

func init() {
	PrivateKeysAnvil[0], _ = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	PrivateKeysAnvil[1], _ = crypto.HexToECDSA("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	PrivateKeysAnvil[2], _ = crypto.HexToECDSA("5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a")
}

func (t *TxBuilder) CallFaucet(addr string, amount *big.Int) (common.Hash, error) {
	client, err := ethclient.Dial(t.endpoint)
	if err != nil {
		return [32]byte{}, err
	}

	publicKey := t.faucetPrivKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return [32]byte{}, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return [32]byte{}, err
	}

	value := amount
	gasLimit := uint64(100000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return [32]byte{}, err
	}

	toAddress := common.HexToAddress(addr)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return [32]byte{}, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), t.faucetPrivKey)
	if err != nil {
		return [32]byte{}, err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return [32]byte{}, err
	}

	logger.LogDebug(fmt.Sprintf("[backend] faucet tx sent with hash: %s", signedTx.Hash().Hex()))

	return signedTx.Hash(), nil
}
