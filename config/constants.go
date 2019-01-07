package config

const (
	// StorageBasePath base path where user data is stored in vault
	// Example: <StorageBasePath>/<user-uuid>
	StorageBasePath = "users/"

	// Entropy is default  length of the bits in the entropy
	Entropy = 256

	// BitsharesDerivationPath used to hard code BTS to some derivation path.
	BitsharesDerivationPath = "m/44'/69'/69'/69/69"
)

// supported log levels
const (
	Info  = "INFO"
	Error = "ERROR"
	Debug = "DEBUG"
	Fatal = "FATAL"
)
