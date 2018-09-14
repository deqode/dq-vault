package adapter

import (
	"fmt"

	"gitlab.com/arout/Vault/lib/adapter/baseadapter"
	"gitlab.com/arout/Vault/lib/bip44coins"
)

// GetAdapter returns suitable adapter depending on coin type
func GetAdapter(coinType uint16, seed []byte, derivationPath string, isDev bool) (baseadapter.IBlockchainAdapter, error) {
	switch coinType {
	case bip44coins.Bitcoin, bip44coins.TestNet:
		return NewBitcoinAdapter(seed, derivationPath, isDev), nil
	case bip44coins.Ether:
		return NewEthereumAdapter(seed, derivationPath, isDev), nil
	}

	return nil, fmt.Errorf("Unable to find suitable adapter.\nUnsupported coin type %v", coinType)
}
