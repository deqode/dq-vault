package api

import (
	"context"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathPassphrase corresponds to POST gen/passphrase.
func (b *backend) pathKeypair(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := validateFields(req, d); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	// obtain entropy length, passphrase, uuid
	entropyLength := d.Get("entropy").(int)
	passphrase := d.Get("passphrase").(string)
	uuid := d.Get("uuid").(string)

	// generated storage path to store user info
	storagePath := "users/" + uuid

	// TODO: remove current node implementation
	// and store only desired data
	derivationPath := "m/44'/60'/0'/0/0"

	// Obtain all existing UUID's from DB
	vals, err := req.Storage.List(ctx, "users/")
	if err != nil {
		return nil, err
	}

	// check if UUID exists
	for i := 0; i < len(vals); i++ {
		if uuid == vals[i] {
			return nil, logical.CodedError(http.StatusUnprocessableEntity, "UUID already exists")
		}
	}

	// Check if user provided UUID or not
	// TODO: use empty string as default if not provided
	if uuid == "x#y" {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, "Provide a valid UUID")
	}

	// For valid entropy length check
	if entropyLength < 128 || entropyLength%32 != 0 || entropyLength > 256 {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, "Invalid bip32 entropy length")
	}

	// generate entropy of desired length
	entropy, err := bip39.NewEntropy(entropyLength)
	checkError(err, "Error generating entropy")

	// obtain mnemonics from entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	checkError(err, "Error generating mnemonics")

	if !bip39.IsMnemonicValid(mnemonic) {
		log.Fatalf("Generated mnemonic is not valid")
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic, passphrase)
	checkError(err, "Error creating wallet")

	path := hdwallet.MustParseDerivationPath(derivationPath)
	account, err := wallet.Derive(path, true)
	checkError(err, "Error deriving child node")

	privateKey, err := wallet.PrivateKeyBytes(account)
	privateKeyHex := hexutil.Encode(privateKey)[2:]
	checkError(err, "Error generating privatekey")

	publicKey, err := wallet.PublicKeyBytes(account)
	publicKeyHex := hexutil.Encode(publicKey)
	checkError(err, "Error generating publickey")

	address, err := wallet.AddressBytes(account)
	addressHex := hexutil.Encode(address)
	checkError(err, "Error generating address")

	node := &node{
		PrivateKey: privateKeyHex,
		PublicKey:  publicKeyHex,
		Path:       derivationPath,
		Address:    addressHex,
	}

	store, err := logical.StorageEntryJSON(storagePath, node)
	checkError(err, "Error storing keys in vault")
	checkError(req.Storage.Put(ctx, store), "")

	return &logical.Response{
		Data: map[string]interface{}{
			"mnemonic":      mnemonic,
			"privateKeyHex": privateKeyHex,
			"publicKeyHex":  publicKeyHex,
			"address":       addressHex,
			"path":          derivationPath,
		},
	}, nil
}
