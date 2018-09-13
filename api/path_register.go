package api

import (
	"context"
	"fmt"
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
		b.Backend.Logger().Info(fmt.Sprintf("\n[ERROR ] register: %v", err.Error()))
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// obatin username
	username := d.Get("username").(string)

	// obtain mnemonic and passphrase of user
	mnemonic := d.Get("mnemonic").(string)
	passphrase := d.Get("passphrase").(string)

	// default entropy length
	entropyLength := config.Entropy

	b.Backend.Logger().Info(fmt.Sprintf("\n[INFO ] register: request username=%v ", username))

	// generate new random UUID
	uuid := helpers.NewUUID()
	for helpers.UUIDExists(ctx, req, uuid) {
		uuid = helpers.NewUUID()
	}

	// generated storage path to store user info
	storagePath := config.StorageBasePath + uuid

	if mnemonic == "" {
		// generate new mnemonics if not provided by user
		// obtain mnemonics from entropy
		mnemonic, err = lib.MnemonicFromEntropy(entropyLength)
		if err != nil {
			b.Backend.Logger().Info(fmt.Sprintf("\n[ERROR ] register: %v", err.Error()))
			return nil, logical.CodedError(http.StatusExpectationFailed, err.Error())
		}
	}

	// check if mnemonic is valid or not
	if !lib.IsMnemonicValid(mnemonic) {
		b.Backend.Logger().Info(fmt.Sprintf("\n[ERROR ] register: invalid mnemonic=[%v]", mnemonic))
		return nil, logical.CodedError(http.StatusExpectationFailed, "Invalid Mnemonic")
	}

	// create object to store user information
	user := &helpers.User{
		Username:   username,
		UUID:       uuid,
		Mnemonic:   mnemonic,
		Passphrase: passphrase,
	}

	// creates strorage entry with user JSON encoded value
	store, err := logical.StorageEntryJSON(storagePath, user)
	if err != nil {
		b.Backend.Logger().Info(fmt.Sprintf("\n[ERROR ] register: %v", err.Error()))
		return nil, logical.CodedError(http.StatusExpectationFailed, err.Error())
	}

	// put user information in store
	if err = req.Storage.Put(ctx, store); err != nil {
		b.Backend.Logger().Info(fmt.Sprintf("\n[ERROR ] register: %v", err.Error()))
		return nil, logical.CodedError(http.StatusExpectationFailed, err.Error())
	}

	b.Backend.Logger().Info(fmt.Sprintf("\n[INFO ] register: user registered uuid=%v username=%v", uuid, username))

	// return response
	return &logical.Response{
		Data: map[string]interface{}{
			"username":   username,
			"uuid":       uuid,
			"mnemonic":   mnemonic,
			"passphrase": passphrase,
		},
	}, nil
}
