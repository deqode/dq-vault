package lib

import (
	"encoding/json"
	"fmt"
	"reflect"

	"gitlab.com/arout/Vault/lib/bip44coins"
)

// DecodeRawTransaction decodes input payload into suitable raw transaction
// depending on cointype provided
// returns error if no supported coin type found
// TODO: improve this
func DecodeRawTransaction(coinType uint16, payload string) (IRawTx, error) {
	var err error
	switch coinType {
	case bip44coins.Bitcoin, bip44coins.TestNet:
		var tx BitcoinRawTx
		if err = json.Unmarshal([]byte(payload), &tx); err != nil ||
			reflect.DeepEqual(tx, BitcoinRawTx{}) {
			return tx, fmt.Errorf("Unable to decode payload=[%v] into coin type %v", payload, coinType)
		}
		return tx, err

	case bip44coins.Ether:
		var tx EthereumRawTx
		if err = json.Unmarshal([]byte(payload), &tx); err != nil ||
			reflect.DeepEqual(tx, EthereumRawTx{}) {
			return tx, fmt.Errorf("Unable to decode payload=[%v] into coin type %v", payload, coinType)
		}
		return tx, err
	case bip44coins.Bitshares:
		var tx BitsharesRawTx
		if err = json.Unmarshal([]byte(payload), &tx); err != nil ||
			reflect.DeepEqual(tx, BitsharesRawTx{}) {
			return tx, fmt.Errorf("Unable to decode payload=[%v] into coin type %v", payload, coinType)
		}
		return tx, err
	}

	return EthereumRawTx{}, fmt.Errorf("Unsupported coin type %v", coinType)
}
