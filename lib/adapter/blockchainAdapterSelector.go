package adapter

import (
	"gitlab.com/arout/Vault/lib/adapter/baseadapter"
	"gitlab.com/arout/Vault/lib/bip44coins"
)

func GetAdapter(coinType uint16, seed []byte, derivationPath string) baseadapter.IBlockchainAdapter {
	switch coinType {
	case bip44coins.Ether:
		return NewEthereumAdapter(seed, derivationPath)
	}

	return nil
}
