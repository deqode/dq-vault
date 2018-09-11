package baseadapter

import "gitlab.com/arout/Vault/lib"

type IBlockchainAdapter interface {
	DerivePrivateKey() (string, error)
	GetBlockchainNetwork() string
	SetEnvironmentToDevelopment()
	SetEnvironmentToProduction()
	CreateSignedTransaction(lib.IRawTx) (string, error)
}

type BlockchainAdapter struct {
	IBlockchainAdapter
	Seed           []byte
	DerivationPath string
	PrivateKey     string
	IsDev          bool
}
