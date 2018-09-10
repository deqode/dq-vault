package lib

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	hdwallet "github.com/go-ethereum-hdwallet"
)

func GetECDSAKeys(seed []byte, derivationPath string) (ECDSAKeyPair, error) {
	// TODO: validate derivationPath
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return ECDSAKeyPair{}, err
	}

	path := hdwallet.MustParseDerivationPath(derivationPath)
	account, err := wallet.Derive(path, true)

	if err != nil {
		return ECDSAKeyPair{}, err
	}

	privateKey, err := wallet.PrivateKeyBytes(account)
	if err != nil {
		return ECDSAKeyPair{}, err
	}

	publicKey, err := wallet.PublicKeyBytes(account)

	if err != nil {
		return ECDSAKeyPair{}, err
	}

	privateKeyHex := hexutil.Encode(privateKey)[2:]
	publicKeyHex := hexutil.Encode(publicKey)

	return ECDSAKeyPair{privateKeyHex, publicKeyHex}, nil
}
