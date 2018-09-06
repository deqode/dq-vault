package gen

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathSignature(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	//get uuid of user
	uuid := d.Get("uuid").(string)

	// TODO: check if UUID is provided or not
	// Handle error cases if UUID is wrong
	path := "users/" + uuid

	// Sample data to create raw transaction
	nonce := uint64(0)
	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice := big.NewInt(30000000000)      // in wei (30 gwei)
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")

	entry, err := req.Storage.Get(ctx, path)

	if err != nil {
		return nil, err
	}

	var node node // create object of the actual struct stored

	err = entry.DecodeJSON(&node)

	//generates ecdsa type key
	privateKey, err := crypto.HexToECDSA(node.PrivateKey)

	if err != nil {
		return nil, err
	}

	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	//sign the transaction
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)

	if err != nil {
		return nil, err
	}

	ts := types.Transactions{signedTx}

	// TODO: provide chainId
	txHex := fmt.Sprintf("%x", ts.GetRlp(0))

	//send signature back to the user
	return &logical.Response{
		Data: map[string]interface{}{
			"signature": txHex,
		},
	}, nil

}
