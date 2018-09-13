package lib

// IRawTx Raw transaction interface
// to enable decoding of all variants of raw transactions (JSON)
type IRawTx interface{}

// EthereumRawTx Ethereum raw transaction implements IRawTx
// to store raw Ethereum JSON payload
type EthereumRawTx struct {
	Nonce    uint64 `json:"nonce"`
	Value    uint64 `json:"value"`
	GasLimit uint64 `json:"gasLimit"`
	GasPrice uint64 `json:"gasPrice"`
	To       string `json:"to"`
	Data     string `json:"data"`
	ChainID  int64  `json:"chainId"`
	IRawTx
}

// UTXO to store bitcoin based transaction's inputs
type UTXO struct {
	TxHash string `json:"txhash"`
	Vout   uint32 `json:"vout"`
}

// PayeeAddress to store bitcoin based
// transactions output details
type PayeeAddress struct {
	Address string `json:"address"`
	Amount  int64  `json:"amount"`
}

// BitcoinRawTx stores bitcoin based raw transaction payloads
// stores input UTXO's and output Addresses
// implements IRawTx
type BitcoinRawTx struct {
	Inputs  []*UTXO         `json:"inputs"`
	Outputs []*PayeeAddress `json:"outputs"`
	IRawTx
}
