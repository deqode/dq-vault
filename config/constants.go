package config

const (
	// StorageBasePath base path where user data is stored in vault
	// Example: <StorageBasePath>/<user-uuid>
	StorageBasePath = "users/"

	// Entropy is default  length of the bits in the entropy
	Entropy = 256
)

// supported log levels
const (
	Info  = "INFO"
	Error = "ERROR"
	Debug = "DEBUG"
	Fatal = "FATAL"
)
