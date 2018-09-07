package api

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"gitlab.com/arout/Vault/api/helpers"
)

func (b *backend) pathRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	var test = "\n"
	// Obtain all existing UUID's from DB
	vals, err := req.Storage.List(ctx, "users/")
	if err != nil {
		return nil, err
	}

	// check if UUID exists
	for i := 0; i < len(vals); i++ {
		path := "users/" + vals[i]
		entry, err := req.Storage.Get(ctx, path)
		helpers.CheckError(err, "")

		// create object of the actual struct stored
		var userInfo helpers.User
		err = entry.DecodeJSON(&userInfo)
		helpers.CheckError(err, "")

		test += userInfo.UUID + ",	" + userInfo.Mnemonic + ",	" + userInfo.Passphrase + "\n"
	}

	//send signature back to the user
	return &logical.Response{
		Data: map[string]interface{}{
			"data": test,
		},
	}, nil
}
