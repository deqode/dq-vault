package api

import (
	"context"
	"net/http"

	"gitlab.com/arout/Vault/lib"
	"gitlab.com/arout/Vault/lib/adapter"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"gitlab.com/arout/Vault/api/helpers"
)

func (b *backend) pathSignature(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := helpers.ValidateFields(req, d); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// obatin data provided
	uuid := d.Get("uuid").(string)
	derivationPath := d.Get("path").(string)
	coinType := d.Get("coinType").(int)
	payload := d.Get("payload").(string)

	// validate data provided
	if err := helpers.ValidateData(ctx, req, uuid, derivationPath); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// TODO: validate data
	rawTransaction, err := lib.DecodeRawTransaction(uint16(coinType), payload)
	if err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// path where user data is stored in vault
	path := "users/" + uuid
	entry, err := req.Storage.Get(ctx, path)
	helpers.CheckError(err, "")

	// create object of the actual struct stored
	var userInfo helpers.User
	err = entry.DecodeJSON(&userInfo)
	helpers.CheckError(err, "")

	// obtain seed from mnemonic and passphrase
	seed, err := lib.SeedFromMnemonic(userInfo.Mnemonic, userInfo.Passphrase)

	// get adapter based on cointype
	adapter := adapter.GetAdapter(60, seed, derivationPath)

	// generate ECDSA keys
	privateKey, err := adapter.GetKeyPair()

	// sign payload sent by application server
	txHex, err := adapter.CreateSignedTransaction(rawTransaction)
	helpers.CheckError(err, "")

	//send signature back to the user
	return &logical.Response{
		Data: map[string]interface{}{
			"raw":       rawTransaction,
			"private":   privateKey,
			"signature": txHex,
		},
	}, nil
}
