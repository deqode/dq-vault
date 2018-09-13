package lib

import (
	"fmt"
	"reflect"
	"testing"
)

type input struct {
	coinType uint16
	payload  string
}

type output struct {
	result IRawTx
	err    string
}
type testpair struct {
	input  input
	output output
}

var decoderTests = []testpair{
	{
		input: input{
			coinType: 60,
			payload:  "{\"nonce\":0,\"value\":100000000000000000,\"gasLimit\":21000,\"gasPrice\":5,\"to\":\"0x95273d64876408E0eDa01a45775Efc2Df6d1CfaC\",\"data\":\"\",\"chainId\":3}",
		}, output: output{
			result: EthereumRawTx{
				Nonce:    0,
				Value:    100000000000000000,
				GasLimit: 21000,
				GasPrice: 5,
				To:       "0x95273d64876408E0eDa01a45775Efc2Df6d1CfaC",
				Data:     "",
				ChainID:  3,
			}, err: "",
		},
	}, {
		input: input{
			coinType: 0,
			payload:  "{\"inputs\":[{\"txhash\":\"b31695ff693b196d41600266d82bdf1092a4a55be608f41e1bde985408b16774\",\"vout\":0}],\"outputs\":[{\"address\":\"3BGgKxAsqoFyouTgUJGW3TAJdvYrk43Jr5\",\"amount\":91234}]}",
		}, output: output{
			result: BitcoinRawTx{
				Inputs:  []*UTXO{&UTXO{TxHash: "b31695ff693b196d41600266d82bdf1092a4a55be608f41e1bde985408b16774", Vout: 0}},
				Outputs: []*PayeeAddress{&PayeeAddress{Address: "3BGgKxAsqoFyouTgUJGW3TAJdvYrk43Jr5", Amount: 91234}},
			},
			err: "",
		},
	}, {
		input: input{
			coinType: 0,
			payload:  "{\"nonce\":0,\"value\":100000000000000000,\"gasLimit\":21000,\"gasPrice\":5,\"to\":\"0x95273d64876408E0eDa01a45775Efc2Df6d1CfaC\",\"data\":\"\",\"chainId\":3}",
		}, output: output{
			result: BitcoinRawTx{},
			err:    fmt.Sprintf("Unable to decode payload=[%v] into coin type %v", "{\"nonce\":0,\"value\":100000000000000000,\"gasLimit\":21000,\"gasPrice\":5,\"to\":\"0x95273d64876408E0eDa01a45775Efc2Df6d1CfaC\",\"data\":\"\",\"chainId\":3}", 0),
		},
	},
}

func TestDecoder(t *testing.T) {
	for _, pair := range decoderTests {
		res, err := DecodeRawTransaction(pair.input.coinType, pair.input.payload)
		if !reflect.DeepEqual(res, pair.output.result) ||
			(err != nil && err.Error() != pair.output.err) {
			t.Error(
				"For", pair.input.payload,
				"\nExpected", pair.output.result,
				"\nGot", res, err,
				"\nErr", err,
			)
		}
	}
}
