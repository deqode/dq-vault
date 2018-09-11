package baseadapter

import "gitlab.com/arout/Vault/lib"

type BitcoinBaseAdapter struct {
	BlockchainAdapter
}

func (b *BitcoinBaseAdapter) DerivePrivateKey() (string, error) { return "", nil }

func (b *BitcoinBaseAdapter) GetWalletAddress() {}

func (b *BitcoinBaseAdapter) GetBlockchainNetwork() string {
	if b.IsDev {
		return "testnet"
	}
	return "mainnet"
}

func (b *BitcoinBaseAdapter) SetEnvironmentToDevelopment() {
	b.IsDev = true
}

func (b *BitcoinBaseAdapter) SetEnvironmentToProduction() {
	b.IsDev = false
}

func (b *BitcoinBaseAdapter) CreateSignedTransaction(payload lib.IRawTx) (string, error) {
	return "", nil
}
