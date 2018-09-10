package lib

import (
	"encoding/json"
	"fmt"

	"gitlab.com/arout/Vault/lib/bip44coins"
)

type IRawTx interface{}

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

func DecodeRawTransaction(coinType uint16, payload string) (IRawTx, error) {
	switch coinType {
	case bip44coins.Ether:
		var tx EthereumRawTx
		err := json.Unmarshal([]byte(payload), &tx)
		return tx, err
	}

	return EthereumRawTx{}, fmt.Errorf("Unsupported coin type %v", coinType)
}
