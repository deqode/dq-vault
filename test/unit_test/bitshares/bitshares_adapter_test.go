package tests

import (
	"testing"

	log "github.com/mgutz/logxi/v1"
	"github.com/deqode/dq-vault/lib/adapter"
)

var logger = log.New("tests")

type btsAdapterCommonInput struct {
	seed           []byte
	derivationPath string
	isDev          bool
}

type btsAdapterPrivateKeyPair struct {
	input      btsAdapterCommonInput
	privateKey string
}

type btsAdapterPublicKeyPair struct {
	input     btsAdapterCommonInput
	publicKey string
}

type btsAdapterAddressPair struct {
	input   btsAdapterCommonInput
	address string
}

type btsAdapterSignaturePair struct {
	input     btsAdapterCommonInput
	payload   string
	signature string
}

var btsAdapterPrivateKeyTests = []btsAdapterPrivateKeyPair{
	{
		btsAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/69'/69'/69/69", false},
		"5JoEZ5hMJgUx57RMgTt8oxCKKdp9zrW3ZHFegCBNPrG7wgM9bFv",
	}, {
		btsAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/69'/69'/69/69", true},
		"5JoEZ5hMJgUx57RMgTt8oxCKKdp9zrW3ZHFegCBNPrG7wgM9bFv",
	}, {
		btsAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/69'/69'/69/69", false},
		"5JdACKsE97FNu9KKmTzzmyrhPnJjjKv9QDLYnGotJYviVazsjc5",
	}, {
		btsAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/0'/0'/0/1", false},
		"5JdACKsE97FNu9KKmTzzmyrhPnJjjKv9QDLYnGotJYviVazsjc5",
	}, {
		btsAdapterCommonInput{[]byte{225, 124, 247, 221, 250, 142, 182, 27, 152, 9, 191, 103, 209, 22, 144, 103, 190, 21, 157, 44, 240, 157, 1, 164, 162, 19, 87, 86, 38, 60, 13, 87, 184, 37, 27, 240, 171, 41, 33, 2, 250, 173, 64, 47, 93, 228, 13, 240, 96, 4, 156, 74, 21, 5, 184, 71, 163, 191, 74, 51, 48, 113, 87, 132}, "m/44'/69'/69'/69/69", false},
		"5K44658mCyJQtmT2TtwgDJtCmduDGibq64mQ9MBBQVWtWSMLvYu",
	},
}

var btsAdapterPublicKeyTests = []btsAdapterPublicKeyPair{
	{
		btsAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/69'/69'/69/69", false},
		"BTS6ZSvDjmeKj6QfXkYJpX3DKcxzKzggpzQbNaxC3XhtamQNqUBrc",
	}, {
		btsAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/69'/69'/69/69", true},
		"BTS6ZSvDjmeKj6QfXkYJpX3DKcxzKzggpzQbNaxC3XhtamQNqUBrc",
	}, {
		btsAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/69'/69'/69/69", false},
		"BTS6satGrftpQBrKRuxgoV2pBiHismf2XEBmqbYXBcVNxq4a9bMFC",
	}, {
		btsAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/0'/0'/0/1", false},
		"BTS6satGrftpQBrKRuxgoV2pBiHismf2XEBmqbYXBcVNxq4a9bMFC",
	}, {
		btsAdapterCommonInput{[]byte{225, 124, 247, 221, 250, 142, 182, 27, 152, 9, 191, 103, 209, 22, 144, 103, 190, 21, 157, 44, 240, 157, 1, 164, 162, 19, 87, 86, 38, 60, 13, 87, 184, 37, 27, 240, 171, 41, 33, 2, 250, 173, 64, 47, 93, 228, 13, 240, 96, 4, 156, 74, 21, 5, 184, 71, 163, 191, 74, 51, 48, 113, 87, 132}, "m/44'/69'/69'/69/69", false},
		"BTS8kJJjEV27b6MZW9n6toS3sTCE8mw9wULnF7J5PvKBn1uDuViFH",
	},
}

var btsAdapterAddressTests = []btsAdapterAddressPair{
	{
		btsAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/69'/69'/69/69", false},
		"BTSJufuQ89rGcSJf9Ko8JzgNm4x8b9qs6w57",
	}, {
		btsAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/69'/69'/69/69", true},
		"BTSJufuQ89rGcSJf9Ko8JzgNm4x8b9qs6w57",
	}, {
		btsAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/69'/69'/69/69", false},
		"BTSHfSGqhP7uMrNm8uMw9CU5DdLitg7qDs6b",
	}, {
		btsAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/0'/0'/0/1", false},
		"BTSHfSGqhP7uMrNm8uMw9CU5DdLitg7qDs6b",
	}, {
		btsAdapterCommonInput{[]byte{225, 124, 247, 221, 250, 142, 182, 27, 152, 9, 191, 103, 209, 22, 144, 103, 190, 21, 157, 44, 240, 157, 1, 164, 162, 19, 87, 86, 38, 60, 13, 87, 184, 37, 27, 240, 171, 41, 33, 2, 250, 173, 64, 47, 93, 228, 13, 240, 96, 4, 156, 74, 21, 5, 184, 71, 163, 191, 74, 51, 48, 113, 87, 132}, "m/44'/69'/69'/69/69", false},
		"BTSJ9Bv8DN29vCGod6Z7NEMQ1kHAMFyK832T",
	},
}

var btsAdapterSignatureTests = []btsAdapterSignaturePair{
	{
		btsAdapterCommonInput{[]byte{231, 123, 253, 32, 91, 177, 57, 41, 83, 210, 141, 254, 70, 155, 155, 209, 146, 239, 121, 33, 115, 236, 10, 103, 15, 213, 59, 14, 171, 113, 96, 133, 224, 169, 71, 197, 252, 254, 148, 145, 81, 131, 232, 109, 136, 38, 30, 96, 164, 19, 58, 42, 207, 81, 15, 139, 154, 151, 104, 139, 171, 132, 186, 122}, "m/44'/69'/69'/69/69", false},
		"{\"transactionDigest\":\"3aef3997194701308d57a65214a7a015d98382ab66a9bc0d655de80842b6bfc5aede09dd6e161ca9095c0105d1d8070000000000001111050007616e6b69743131010000000001021500e918e7ca8c63e40472c9a2ab28665d06a41e78d034ee1b2ff2b3635d02e5010000010000000001021500e918e7ca8c63e40472c9a2ab28665d06a41e78d034ee1b2ff2b3635d02e5010000021500e918e7ca8c63e40472c9a2ab28665d06a41e78d034ee1b2ff2b3635d02e5050000000000000000\"}",
		"2040d87e6b3b4f87debbd8d393852930ad837a391bbea3f9794de76271d9c6ddb049b82700ea23ebb55eb5e8eabbb9e107a95d3f6440ac0458669a25cbb2a9b78f",
	},
}

func TestBTSPrivateKey(t *testing.T) {
	for _, pair := range btsAdapterPrivateKeyTests {
		adapter := adapter.NewBitsharesAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

		wif, _ := adapter.DerivePrivateKey(logger)
		if wif != pair.privateKey {
			t.Error(
				"Seed", pair.input.seed,
				"\nPath", pair.input.derivationPath,
				"\nExpected WIF", pair.privateKey,
				"\nGot WIF", wif,
			)
		}
	}
}

func TestBTSPublicKey(t *testing.T) {
	for _, pair := range btsAdapterPublicKeyTests {
		adapter := adapter.NewBitsharesAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

		publicKey, _ := adapter.DerivePublicKey(logger)
		if publicKey != pair.publicKey {
			t.Error(
				"Seed", pair.input.seed,
				"\nPath", pair.input.derivationPath,
				"\nExpected PublicKey", pair.publicKey,
				"\nGot PublicKey", publicKey,
			)
		}
	}
}

func TestBTSAddress(t *testing.T) {
	for _, pair := range btsAdapterAddressTests {
		adapter := adapter.NewBitsharesAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

		address, _ := adapter.DeriveAddress(logger)
		if address != pair.address {
			t.Error(
				"Seed", pair.input.seed,
				"\nPath", pair.input.derivationPath,
				"\nExpected Address", pair.address,
				"\nGot Address", address,
			)
		}
	}
}

func TestBTSSignature(t *testing.T) {
	for _, pair := range btsAdapterSignatureTests {
		adapter := adapter.NewBitsharesAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

		adapter.DerivePrivateKey(logger)
		signature, _ := adapter.CreateSignedTransaction(pair.payload, logger)
		if signature != pair.signature {
			t.Error(
				"Seed", pair.input.seed,
				"\nPath", pair.input.derivationPath,
				"\nExpected Signature", pair.signature,
				"\nGot Signature", signature,
			)
		}
	}
}
