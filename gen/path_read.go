package gen

import (
	"context"
	"math/big"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (b *backend) readKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	path := "store" //change path accordingly

	entry, err := req.Storage.Get(ctx, path)

	if err != nil {
		return nil, err
	}

	var key Key // create object of the actual struct stored

	err = entry.DecodeJSON(&key)

	//generates ecdsa type key 
	privateKey, err := crypto.HexToECDSA(key.PrivateKey) 

	if err != nil {
		return nil, err
	}

	//dummy transaction
	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000) // in units
	gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")

	var data []byte
	tx := types.NewTransaction(0, toAddress, value, gasLimit, gasPrice, data)

	//sign the transaction
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)

	if err != nil {
	 return nil, err
	}

	//send signature back to the user
	return &logical.Response{
		Data: map[string]interface{}{	
			"signature": signedTx.Hash().Hex(),
		},
	}, nil 
}