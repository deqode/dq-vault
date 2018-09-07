package helpers

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// CreateTransaction - creates transaction from provided data
func CreateTransaction(privateKeyHex string) (string, error) {
	//generates ecdsa type key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	CheckError(err, "")

	tx := dummyTx()

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(3)), privateKey)
	CheckError(err, "")

	ts := types.Transactions{signedTx}
	txHex := fmt.Sprintf("%x", ts.GetRlp(0))

	return txHex, nil
}

func dummyTx() *types.Transaction {
	// Sample data to create raw transaction
	nonce := uint64(0)
	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice := big.NewInt(30000000000)      // in wei (30 gwei)
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")

	var data []byte
	return types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
}
