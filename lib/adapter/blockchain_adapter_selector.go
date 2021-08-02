package adapter

import (
	"fmt"

	"vault/lib/adapter/baseadapter"
	"vault/lib/bip44coins"
)

// GetAdapter returns suitable adapter depending on coin type
func GetAdapter(coinType uint16, seed []byte, derivationPath string) (baseadapter.IBlockchainAdapter, error) {
	switch coinType {
	case bip44coins.Bitcoin:
		return NewBitcoinAdapter(seed, derivationPath, false), nil
	case bip44coins.TestNet:
		return NewBitcoinAdapter(seed, derivationPath, true), nil
	case bip44coins.Ether:
		return NewEthereumAdapter(seed, derivationPath, false), nil
	case bip44coins.Bitshares:
		return NewBitsharesAdapter(seed, derivationPath, false), nil
	}

	return nil, fmt.Errorf("Unable to find suitable adapter.\nUnsupported coin type %v, bip44coins.Bitshares= %v bip44coins.Bitshares==coinType= %v", coinType, bip44coins.Bitshares, bip44coins.Bitshares == coinType)
}
