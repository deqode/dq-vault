package api

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/arout/Vault/api/helpers"
	"gitlab.com/arout/Vault/config"
	"gitlab.com/arout/Vault/lib"
	"gitlab.com/arout/Vault/lib/adapter"
	"gitlab.com/arout/Vault/logger"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathAddress(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	backendLogger := b.Backend.Logger()
	if err := helpers.ValidateFields(req, d); err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// UUID of user which want to sign transaction
	uuid := d.Get("uuid").(string)

	// derivation path
	derivationPath := d.Get("path").(string)

	// coin type of transaction
	// see supported coinTypes lib/bipp44coins
	coinType := d.Get("coinType").(int)

	logger.Log(backendLogger, config.Info, "address:", fmt.Sprintf("request uuid=%v path=[%v] cointype=%v ", uuid, derivationPath, coinType))

	// validate data provided
	if err := helpers.ValidateData(ctx, req, uuid, derivationPath); err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// path where user data is stored in vault
	path := config.StorageBasePath + uuid
	entry, err := req.Storage.Get(ctx, path)
	if err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// obtain mnemonic, passphrase of user
	var userInfo helpers.User
	err = entry.DecodeJSON(&userInfo)
	if err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// obtain seed from mnemonic and passphrase
	seed, err := lib.SeedFromMnemonic(userInfo.Mnemonic, userInfo.Passphrase)

	// obtains blockchain adapater based on coinType
	adapter, err := adapter.GetAdapter(uint16(coinType), seed, derivationPath)
	if err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// Generates and stores ECDSA private key in adapter
	_, err = adapter.DerivePrivateKey(backendLogger)
	if err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	pubKey, err := adapter.DerivePublicKey(backendLogger)
	if err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	address, err := adapter.DeriveAddress(backendLogger)
	if err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	logger.Log(backendLogger, config.Info, "address:", fmt.Sprintf("\n[INFO ] address: uuid=%v derived publicKey=[%v], address=[%v]", uuid, pubKey, address))

	// Returns publicKey, address as output
	return &logical.Response{
		Data: map[string]interface{}{
			"uuid":      uuid,
			"publicKey": pubKey,
			"address":   address,
		},
	}, nil
}
