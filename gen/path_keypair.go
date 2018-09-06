package gen

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

	// obtain entropy length
	entropyLength := d.Get("entropy").(int)

	//passphrase := d.Get("passphrase").(string) not used, refer line 43

	//get uuid of user
	uid := d.Get("uid").(string)

	//get list of all existing uuids
	vals, err := req.Storage.List(ctx, "users/")
	if err != nil {
		return nil, err
	}

	//improve this if possible, checks with existing uuids
	for i:=0; i<len(vals); i++ {
		if(uid == vals[i]) {
			return nil, logical.CodedError(http.StatusUnprocessableEntity, "uid already exists")
		}
	}

	//checks if value is provided or not
	if uid=="x#y" {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, "provide an user id") 
	}

	storagePath := "users/" + uid
	derivationPath := "m/44'/60'/0'/0/0"

	if entropyLength < 128 || entropyLength%32 != 0 || entropyLength > 256 {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, "Invalid bip32 entropy length")
	}

	// generate entropy of desired length
	entropy, err := bip39.NewEntropy(entropyLength)
	checkError(err, "Error generating entropy")

	// obtain mnemonics from entropy
	mnemonic, err := bip39.NewMnemonic(entropy /*,passphrase*/) //showing error if passphrase is provided
	checkError(err, "Error generating mnemonics")

	if !bip39.IsMnemonicValid(mnemonic) {
		log.Fatalf("Generated mnemonic is not valid")
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
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
