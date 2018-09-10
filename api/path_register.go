package api

import (
	"context"
	"log"
	"net/http"

	"gitlab.com/arout/Vault/api/helpers"
	"gitlab.com/arout/Vault/lib"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathPassphrase corresponds to POST gen/passphrase.
func (b *backend) pathRegister(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := helpers.ValidateFields(req, d); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// obatin data provided
	username := d.Get("username").(string)
	mnemonic := d.Get("mnemonic").(string)
	passphrase := d.Get("passphrase").(string)
	entropyLength := 256

	//TODO: check if username exists or not

	// generate UUID
	uuid := helpers.NewUUID()
	var err error

	// generated storage path to store user info
	storagePath := "users/" + uuid

	if mnemonic == "" {
		// generate new mnemonics if not provided by user
		// obtain mnemonics from entropy
		mnemonic, err = lib.MnemonicFromEntropy(entropyLength)
		helpers.CheckError(err, "Error generating mnemonics")
	}

	if !lib.IsMnemonicValid(mnemonic) {
		log.Fatalf("Mnemonic is not valid")
	}

	user := &helpers.User{
		Username:   username,
		UUID:       uuid,
		Mnemonic:   mnemonic,
		Passphrase: passphrase,
	}

	store, err := logical.StorageEntryJSON(storagePath, user)
	helpers.CheckError(err, "Error storing data in vault")
	helpers.CheckError(req.Storage.Put(ctx, store), "")

	return &logical.Response{
		Data: map[string]interface{}{
			"username":   username,
			"uuid":       uuid,
			"mnemonic":   mnemonic,
			"passphrase": passphrase,
		},
	}, nil
}
