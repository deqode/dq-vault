package adapter

import "gitlab.com/arout/Vault/lib/adapter/baseadapter"

type BitcoinAdapter struct {
	baseadapter.BitcoinBaseAdapter
}

func NewBitcoinAdapter(seed []byte, derivationPath string) *BitcoinAdapter {
	adapter := new(BitcoinAdapter)
	adapter.Seed = seed
	adapter.DerivationPath = derivationPath
	adapter.IsDev = false

	return adapter
}
