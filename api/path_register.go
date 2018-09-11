package api

import (
	"context"
	"net/http"

	"gitlab.com/arout/Vault/api/helpers"
	"gitlab.com/arout/Vault/config"
	"gitlab.com/arout/Vault/lib"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathPassphrase corresponds to POST gen/passphrase.
func (b *backend) pathRegister(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var err error
	if err = helpers.ValidateFields(req, d); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// obatin username
	username := d.Get("username").(string)

	// obtain mnemonic and passphrase of user
	mnemonic := d.Get("mnemonic").(string)
	passphrase := d.Get("passphrase").(string)

	// default entropy length
	entropyLength := config.Entropy

	// generate new random UUID
	uuid := helpers.NewUUID()

	// generated storage path to store user info
	storagePath := config.StorageBasePath + uuid

	if mnemonic == "" {
		// generate new mnemonics if not provided by user
		// obtain mnemonics from entropy
		mnemonic, err = lib.MnemonicFromEntropy(entropyLength)
		if err != nil {
			return nil, logical.CodedError(http.StatusExpectationFailed, err.Error())
		}
	}

	if !lib.IsMnemonicValid(mnemonic) {
		return nil, logical.CodedError(http.StatusExpectationFailed, "Invalid Mnemonic")
	}

	user := &helpers.User{
		Username:   username,
		UUID:       uuid,
		Mnemonic:   mnemonic,
		Passphrase: passphrase,
	}

	store, err := logical.StorageEntryJSON(storagePath, user)
	if err != nil {

	}

	if err = req.Storage.Put(ctx, store); err != nil {
		return nil, logical.CodedError(http.StatusExpectationFailed, err.Error())
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"username":   username,
			"uuid":       uuid,
			"mnemonic":   mnemonic,
			"passphrase": passphrase,
		},
	}, nil
}
