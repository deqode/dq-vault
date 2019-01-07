package test

import (
	"bytes"
	"testing"

	"gitlab.com/arout/Vault/lib"
)

type mnemonicTestPair struct {
	mnemonic   string
	passphrase string
	seed       []byte
}

var mnemonicTests = []mnemonicTestPair{
	{
		"crumble brother salt endless few process maple alter panda trim trumpet usual skate ritual clerk theme image stable air emerge demand truth wear interest",
		"",
		[]byte{24, 237, 186, 111, 109, 126, 253, 50, 113, 111, 1, 210, 61, 158, 59, 221, 138, 181, 248, 11, 255, 243, 230, 213, 155, 49, 156, 7, 194, 255, 253, 89, 116, 121, 73, 136, 40, 151, 87, 62, 76, 218, 234, 190, 76, 153, 151, 170, 99, 41, 41, 94, 57, 118, 18, 40, 199, 60, 61, 12, 252, 183, 131, 148},
	}, {
		"salmon arena devote news actor bubble skull smoke foil mango head catalog oil drastic spell suggest flag fitness echo exhaust fetch derive robust loud",
		"abc",
		[]byte{248, 203, 122, 99, 136, 88, 160, 24, 138, 242, 68, 214, 241, 12, 34, 86, 141, 151, 40, 173, 40, 24, 157, 101, 101, 147, 251, 213, 237, 238, 27, 88, 162, 233, 212, 55, 155, 226, 33, 31, 85, 205, 83, 168, 53, 157, 16, 221, 154, 64, 143, 56, 235, 112, 37, 152, 220, 213, 17, 94, 25, 12, 230, 27},
	}, {
		"purity distance bunker negative journey good dumb service tackle marriage second turkey oppose leg castle require essence bleak paddle chapter animal stomach month immune",
		"blockchain",
		[]byte{171, 137, 160, 246, 1, 222, 42, 19, 155, 55, 100, 0, 78, 140, 204, 143, 71, 32, 188, 114, 254, 19, 170, 185, 55, 222, 49, 37, 200, 45, 73, 22, 57, 197, 204, 153, 204, 80, 93, 177, 3, 119, 178, 119, 149, 228, 225, 52, 215, 213, 128, 151, 120, 255, 72, 110, 36, 211, 104, 196, 134, 189, 61, 136},
	},
}

func TestMnemonic(t *testing.T) {
	for _, pair := range mnemonicTests {
		seed, _ := lib.SeedFromMnemonic(pair.mnemonic, pair.passphrase)
		if !bytes.Equal(seed, pair.seed) {
			t.Error(
				"Mnemonic", pair.mnemonic,
				"\nPassphrase", pair.passphrase,
				"\nExpected Seed", pair.seed,
				"\nGot Seed", seed,
			)
		}
	}
}
