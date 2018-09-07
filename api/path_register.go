package api

import (
	"context"
	"log"
	"net/http"

	"github.com/tyler-smith/go-bip39"
	"gitlab.com/arout/Vault/api/helpers"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathPassphrase corresponds to POST gen/passphrase.
func (b *backend) pathRegister(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := helpers.ValidateFields(req, d); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// obatin data provided
	uuid := d.Get("uuid").(string)
	mnemonic := d.Get("mnemonic").(string)
	passphrase := d.Get("passphrase").(string)
	entropyLength := 256

	// generated storage path to store user info
	storagePath := "users/" + uuid

	test := ""

	// Obtain all existing UUID's from DB
	vals, err := req.Storage.List(ctx, "users/")
	if err != nil {
		return nil, err
	}

	// check if UUID exists
	for i := 0; i < len(vals); i++ {
		test += vals[i]
		if uuid == vals[i] {
			return nil, logical.CodedError(http.StatusUnprocessableEntity, "Provided UUID already exists")
		}
	}

	// Check if user provided UUID or not
	if uuid == "" {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, "Provide a valid UUID")
	}

	if mnemonic == "" {
		// generate entropy of desired length
		entropy, err := bip39.NewEntropy(entropyLength)
		helpers.CheckError(err, "Error generating entropy")

		// generate new mnemonics if not provided by user
		// obtain mnemonics from entropy
		mnemonic, err = bip39.NewMnemonic(entropy)
		helpers.CheckError(err, "Error generating mnemonics")
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		log.Fatalf("Mnemonic is not valid")
	}

	user := &helpers.User{
		UUID:       uuid,
		Mnemonic:   mnemonic,
		Passphrase: passphrase,
	}

	store, err := logical.StorageEntryJSON(storagePath, user)
	helpers.CheckError(err, "Error storing data in vault")
	helpers.CheckError(req.Storage.Put(ctx, store), "")

	return &logical.Response{
		Data: map[string]interface{}{
			"test":       test,
			"uuid":       uuid,
			"mnemonic":   mnemonic,
			"passphrase": passphrase,
		},
	}, nil
}
