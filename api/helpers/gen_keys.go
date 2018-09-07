package helpers

import (
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

// GenerateKeys - generates keys using mnemonics, passphrase and deviation path
func GenerateKeys(mnemonic, passphrase, derivationPath string) (privateKeyHex string) {
	// check if mnemonic is valid
	if !bip39.IsMnemonicValid(mnemonic) {
		log.Fatalf("Generated mnemonic is not valid")
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic, passphrase)
	CheckError(err, "Error creating wallet")

	path := hdwallet.MustParseDerivationPath(derivationPath)
	account, err := wallet.Derive(path, true)
	CheckError(err, "Error deriving child node")

	privateKey, err := wallet.PrivateKeyBytes(account)
	privateKeyHex = hexutil.Encode(privateKey)[2:]
	CheckError(err, "Error generating privatekey")

	return
}
