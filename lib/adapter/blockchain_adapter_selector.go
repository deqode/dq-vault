package adapter

import (
	"fmt"

	"gitlab.com/arout/Vault/lib/adapter/baseadapter"
	"gitlab.com/arout/Vault/lib/bip44coins"
)

// GetAdapter returns suitable adapter depending on coin type
func GetAdapter(coinType uint16, seed []byte, derivationPath string) (baseadapter.IBlockchainAdapter, error) {
	switch coinType {
	case bip44coins.Bitcoin:
		return NewBitcoinAdapter(seed, derivationPath), nil
	case bip44coins.TestNet:
		// TODO: improve this
		// sets network to testnet
		testnetAdapter := NewBitcoinAdapter(seed, derivationPath)
		testnetAdapter.SetEnvironmentToDevelopment()
		return testnetAdapter, nil
	case bip44coins.Ether:
		return NewEthereumAdapter(seed, derivationPath), nil
	}

	return nil, fmt.Errorf("Unable to find suitable adapter.\nUnsupported coin type %v", coinType)
}
