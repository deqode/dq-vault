package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/deqode/dq-vault/api/helpers"
	"github.com/deqode/dq-vault/config"
	"github.com/deqode/dq-vault/lib"
	"github.com/deqode/dq-vault/lib/adapter"
	"github.com/deqode/dq-vault/lib/bip44coins"
	"github.com/deqode/dq-vault/logger"
)

func (b *backend) pathAddress(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	backendLogger := b.logger
	if err := helpers.ValidateFields(req, d); err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// UUID of user required to sign transaction
	uuid := d.Get("uuid").(string)

	// derivation path
	derivationPath := d.Get("path").(string)

	// coin type of transaction
	// see supported coinTypes lib/bipp44coins
	coinType := d.Get("coinType").(int)

	if uint16(coinType) == bip44coins.Bitshares {
		derivationPath = config.BitsharesDerivationPath
	}

	logger.Log(backendLogger, config.Info, "address:", fmt.Sprintf("request path=[%v] cointype=%v ", derivationPath, coinType))

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

	// obtain mnemonic and passphrase of user
	var userInfo helpers.User
	err = entry.DecodeJSON(&userInfo)
	if err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// obtain seed from mnemonic and passphrase
	logger.Log(backendLogger, config.Info, "mnemonic", userInfo.Mnemonic, userInfo.Passphrase)
	seed, err := lib.SeedFromMnemonic(userInfo.Mnemonic, userInfo.Passphrase)

	logger.Log(backendLogger, config.Info, "dp", derivationPath)

	// obtains blockchain adapater based on coinType
	adapter, err := adapter.GetAdapter(uint16(coinType), seed, derivationPath)
	if err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// Generates and stores ECDSA private key in adapter
	priv, err := adapter.DerivePrivateKey(backendLogger)
	if err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	logger.Log(backendLogger, config.Info, "priv", priv)

	pubKey, err := adapter.DerivePublicKey(backendLogger)
	if err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error(), "which state", pubKey)
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	address, err := adapter.DeriveAddress(backendLogger)
	if err != nil {
		logger.Log(backendLogger, config.Error, "address:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	logger.Log(backendLogger, config.Info, "address:", fmt.Sprintf("\n[INFO ] address:  derived publicKey=[%v], address=[%v]", pubKey, address))

	// Returns publicKey and address as output
	return &logical.Response{
		Data: map[string]interface{}{
			"uuid":      uuid,
			"publicKey": pubKey,
			"address":   address,
		},
	}, nil
}
