package test

import (
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/deqode/dq-vault/lib"
)

type privateKeyInput struct {
	seed         []byte
	path         string
	isDev        bool
	isCompressed bool
}
type privateKeyTestPair struct {
	input privateKeyInput
	wif   string
}

var privateKeyTestCases = []privateKeyTestPair{
	{
		privateKeyInput{[]byte{24, 237, 186, 111, 109, 126, 253, 50, 113, 111, 1, 210, 61, 158, 59, 221, 138, 181, 248, 11, 255, 243, 230, 213, 155, 49, 156, 7, 194, 255, 253, 89, 116, 121, 73, 136, 40, 151, 87, 62, 76, 218, 234, 190, 76, 153, 151, 170, 99, 41, 41, 94, 57, 118, 18, 40, 199, 60, 61, 12, 252, 183, 131, 148}, "m/44'/0'/0'/0/1", false, true},
		"L45yYN8spZMsadbH87SHmLAYJwPYQGCHfYv6YFzREcs6Xy66nEW1",
	}, {
		privateKeyInput{[]byte{248, 203, 122, 99, 136, 88, 160, 24, 138, 242, 68, 214, 241, 12, 34, 86, 141, 151, 40, 173, 40, 24, 157, 101, 101, 147, 251, 213, 237, 238, 27, 88, 162, 233, 212, 55, 155, 226, 33, 31, 85, 205, 83, 168, 53, 157, 16, 221, 154, 64, 143, 56, 235, 112, 37, 152, 220, 213, 17, 94, 25, 12, 230, 27}, "m/44'/60'/0'/0/0", false, true},
		"KxtZkrrqyTMoim6UjDeTK9wVyhtxPdcMXWCNV6myBVMp1kJmX9Tv",
	},
}

func TestDerivePrivateKey(t *testing.T) {
	for _, pair := range privateKeyTestCases {
		btcPrivKey, _ := lib.DerivePrivateKey(pair.input.seed, pair.input.path, pair.input.isDev)
		if wif, _ := toWIF(btcPrivKey, pair.input.isDev, pair.input.isCompressed); wif != pair.wif {
			t.Error(
				"For seed", pair.input.seed,
				"\nPath", pair.input.path,
				"\nExpected", pair.wif,
				"\nGot", wif,
			)
		}
	}
}

func toWIF(p *btcec.PrivateKey, isDev bool, isCompressed bool) (string, error) {
	network := &chaincfg.MainNetParams
	if isDev {
		network = &chaincfg.TestNet3Params
	}

	privateWIF, err := btcutil.NewWIF(p, network, isCompressed)
	if err != nil {
		return "", err
	}
	return privateWIF.String(), nil
}
