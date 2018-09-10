package baseadapter

import "gitlab.com/arout/Vault/lib"

type IBlockchainAdapter interface {
	GetKeyPair() (string, error)
	GetWalletAddress()
	GetBlockchainNetwork(bool) string
	CreateSignedTransaction(lib.IRawTx) (string, error)
}

type BlockchainAdapter struct {
	IBlockchainAdapter
	Seed           []byte
	DerivationPath string
	PrivateKey     string
	PublicKey      string
	WalletAddress  string
	Balance        string
}

func NewBlockchainAdapter(seed []byte, derivationPath string) *BlockchainAdapter {
	adapter := new(BlockchainAdapter)
	adapter.Seed = seed
	adapter.DerivationPath = derivationPath

	return adapter
}
