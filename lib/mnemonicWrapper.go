package lib

import (
	"errors"

	"github.com/tyler-smith/go-bip39"
)

// GenerateMnemonic will return a string consisting of the mnemonic words for
// the default entropy = 256.
// If the provide entropy is invalid, an error will be returned.
func GenerateMnemonic() (string, error) {
	return MnemonicFromEntropy(256)
}

// MnemonicFromEntropy will return a string consisting of the mnemonic words for
// the given entropy.
func MnemonicFromEntropy(entropyLength int) (string, error) {
	entropy, err := bip39.NewEntropy(entropyLength)
	if err != nil {
		return "", err
	}

	return bip39.NewMnemonic(entropy)
}

// IsMnemonicValid attempts to verify that the provided mnemonic is valid.
// Validity is determined by both the number of words being appropriate,
// and that all the words in the mnemonic are present in the word list.
func IsMnemonicValid(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

// SeedFromMnemonic creates a hashed seed output given a provided string and password.
// No checking is performed to validate that the string provided is a valid mnemonic.
func SeedFromMnemonic(mnemonic, passphrase string) ([]byte, error) {
	if !IsMnemonicValid(mnemonic) {
		return nil, errors.New("Invalid Mnemonic")
	}
	return bip39.NewSeed(mnemonic, passphrase), nil
}
