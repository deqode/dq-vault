package baseadapter

import "gitlab.com/arout/Vault/lib"

// IBlockchainAdapter Blockchain Adapter Interface
// contains common methods for all Blockchain Adapter variants
type IBlockchainAdapter interface {
	DerivePrivateKey() (string, error)
	GetBlockchainNetwork() string
	SetEnvironmentToDevelopment()
	SetEnvironmentToProduction()
	CreateSignedTransaction(lib.IRawTx) (string, error)
}

// BlockchainAdapter contains common fields for
// all Blockchain Adapter variants
type BlockchainAdapter struct {
	Seed           []byte
	DerivationPath string
	PrivateKey     string
	IsDev          bool
	IBlockchainAdapter
}
