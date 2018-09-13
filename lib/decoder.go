package lib

import (
	"encoding/json"
	"fmt"

	"gitlab.com/arout/Vault/lib/bip44coins"
)

// DecodeRawTransaction decodes input payload into suitable raw transaction
// depending on cointype provided
// returns error if no supported coin type found
func DecodeRawTransaction(coinType uint16, payload string) (IRawTx, error) {
	switch coinType {
	case bip44coins.Bitcoin:
		var tx BitcoinRawTx
		err := json.Unmarshal([]byte(payload), &tx)
		return tx, err
	case bip44coins.Ether:
		var tx EthereumRawTx
		err := json.Unmarshal([]byte(payload), &tx)
		return tx, err
	}

	return EthereumRawTx{}, fmt.Errorf("Unsupported coin type %v", coinType)
}
