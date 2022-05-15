package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deqode/dq-vault/api/helpers"
	"github.com/deqode/dq-vault/config"
	"github.com/deqode/dq-vault/lib"
	"github.com/deqode/dq-vault/lib/adapter"
	"github.com/deqode/dq-vault/lib/bip44coins"
	"github.com/deqode/dq-vault/logger"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathSign(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	backendLogger := b.logger
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

	// data in string hex
	// depends on type of transaction
	payload := d.Get("payload").(string)

	if uint16(coinType) == bip44coins.Bitshares {
		derivationPath = config.BitsharesDerivationPath
	}

	logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("request  path=[%v] cointype=%v payload=[%v]", derivationPath, coinType, payload))

	// validate data provided
	if err := helpers.ValidateData(ctx, req, uuid, derivationPath); err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// path where user data is stored in vault
	path := config.StorageBasePath + uuid
	entry, err := req.Storage.Get(ctx, path)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// obtain mnemonic, passphrase of user
	var userInfo helpers.User
	err = entry.DecodeJSON(&userInfo)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// obtain seed from mnemonic and passphrase
	seed, err := lib.SeedFromMnemonic(userInfo.Mnemonic, userInfo.Passphrase)

	// obtains blockchain adapater based on coinType
	adapter, err := adapter.GetAdapter(uint16(coinType), seed, derivationPath)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// Generates and stores ECDSA private key in adapter
	_, err = adapter.DerivePrivateKey(backendLogger)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// creates signature from raw transaction payload
	txHex, err := adapter.CreateSignature(payload, backendLogger)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("\n[INFO ] signature: created signature signature=[%v]", txHex))

	// Returns signature as output
	return &logical.Response{
		Data: map[string]interface{}{
			"signature": txHex,
		},
	}, nil
}
