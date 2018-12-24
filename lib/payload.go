package lib

import "math/big"

// IRawTx Raw transaction interface
// to enable decoding of all variants of raw transactions (JSON)
type IRawTx interface{}

// EthereumRawTx Ethereum raw transaction implements IRawTx
// to store raw Ethereum JSON payload
type EthereumRawTx struct {
	Nonce    uint64   `json:"nonce"`
	Value    *big.Int `json:"value"`
	GasLimit uint64   `json:"gasLimit"`
	GasPrice *big.Int `json:"gasPrice"`
	To       string   `json:"to"`
	Data     string   `json:"data"`
	ChainID  *big.Int `json:"chainId"`
	IRawTx
}

// BitcoinRawTx stores bitcoin based raw transaction payloads
// stores input UTXO's and output Addresses
// implements IRawTx
type BitcoinRawTx struct {
	Inputs []struct {
		Txhash string `json:"txhash"`
		Vout   uint32 `json:"vout"`
	} `json:"inputs"`
	Outputs []struct {
		Address string `json:"address"`
		Amount  int64  `json:"amount"`
	} `json:"outputs"`
	IRawTx
}

// BitsharesRawTx Ethereum raw transaction implements IRawTx
// to store raw Bitshares JSON payload
// Bitshares' Payload will only have a hex encoded string that would represent the transaction digest.
type BitsharesRawTx struct {
	TransactionDigest string `json:"transactionDigest"`
	IRawTx
}
