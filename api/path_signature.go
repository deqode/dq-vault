package api

import (
	"context"
	"net/http"

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

	// validate data provided
	if err := helpers.ValidateUser(ctx, req, uuid, derivationPath); err != nil {
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

	privateKeyHex := helpers.GenerateKeys(userInfo.Mnemonic, userInfo.Passphrase, derivationPath)

	txHex, err := helpers.CreateTransaction(privateKeyHex)
	helpers.CheckError(err, "")

	//send signature back to the user
	return &logical.Response{
		Data: map[string]interface{}{
			"signature": txHex,
		},
	}, nil
}
