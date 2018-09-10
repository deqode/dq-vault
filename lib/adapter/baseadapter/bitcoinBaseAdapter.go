package baseadapter

type BitcoinBaseAdapter struct {
	BlockchainAdapter
	wif string
}

func (b *BitcoinBaseAdapter) getKeyPair() {}

func (b *BitcoinBaseAdapter) getWalletAddress() {}

func (b *BitcoinBaseAdapter) getBlockchainNetwork(isDev bool) string { return "" }

func (b *BitcoinBaseAdapter) createSignedTransaction() {}
