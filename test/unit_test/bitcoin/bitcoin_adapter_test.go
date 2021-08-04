package tests

import (
	"testing"

	"github.com/deqode/dq-vault/lib/adapter"
	log "github.com/mgutz/logxi/v1"
)

var logger = log.New("tests")

type btcAdapterCommonInput struct {
	seed           []byte
	derivationPath string
	isDev          bool
}

type btcAdapterPrivateKeyPair struct {
	input      btcAdapterCommonInput
	privateKey string
}

type btcAdapterPublicKeyPair struct {
	input     btcAdapterCommonInput
	publicKey string
}

type btcAdapterAddressPair struct {
	input   btcAdapterCommonInput
	address string
}

type btcAdapterSignaturePair struct {
	input     btcAdapterCommonInput
	payload   string
	signature string
}

var btcAdapterPrivateKeyTests = []btcAdapterPrivateKeyPair{
	{
		btcAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/0'/0'/0/0", false},
		"L1UYhwzbSSMWt55t29Kg6XbuXKP8V1jttq9V2cnkxFSmoRbgdF67",
	}, {
		btcAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/1'/0'/0/0", true},
		"cUPh58DMqiN61PGtriRUUHYJBdW4NK74Eo3ARyjH6siSGZvKrZWm",
	}, {
		btcAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/0'/0'/0/0", false},
		"L2D5HF27c93QPKLtfutrqC9ruXMsjRqoQiByzwTLNbsSnhX544Es",
	}, {
		btcAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/0'/0'/0/1", false},
		"L26p4HPCcUxqbvv3DvWp92Zv9tHQnptLmkatwRL3otUnEd38tXx3",
	}, {
		btcAdapterCommonInput{[]byte{225, 124, 247, 221, 250, 142, 182, 27, 152, 9, 191, 103, 209, 22, 144, 103, 190, 21, 157, 44, 240, 157, 1, 164, 162, 19, 87, 86, 38, 60, 13, 87, 184, 37, 27, 240, 171, 41, 33, 2, 250, 173, 64, 47, 93, 228, 13, 240, 96, 4, 156, 74, 21, 5, 184, 71, 163, 191, 74, 51, 48, 113, 87, 132}, "m/44'/0'/0'/0/0", false},
		"KzpGssz5Ey3c5CyesSnFU31UMnUbgEQaNouoQBehxneXth8PiWxv",
	},
}

var btcAdapterPublicKeyTests = []btcAdapterPublicKeyPair{
	{
		btcAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/0'/0'/0/0", false},
		"027bc1528f2498d41981d867ab97833e99120dcd30b7e2d36f703bb8c7fbb2d8f3",
	}, {
		btcAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/1'/0'/0/0", true},
		"03fb9026b85efe189af4995f22dfa822c20d359eb54e9053c2049a4462824df503",
	}, {
		btcAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/0'/0'/0/0", false},
		"02eef41ef0e33e68d86ed4e0ac5d4a2ee3f5230fdb67166ff27f324189aa7cb380",
	}, {
		btcAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/0'/0'/0/1", false},
		"021d3075c181d2a22e8c77346512ecc522c71f691dafc6dbe54cb17060556fc247",
	}, {
		btcAdapterCommonInput{[]byte{225, 124, 247, 221, 250, 142, 182, 27, 152, 9, 191, 103, 209, 22, 144, 103, 190, 21, 157, 44, 240, 157, 1, 164, 162, 19, 87, 86, 38, 60, 13, 87, 184, 37, 27, 240, 171, 41, 33, 2, 250, 173, 64, 47, 93, 228, 13, 240, 96, 4, 156, 74, 21, 5, 184, 71, 163, 191, 74, 51, 48, 113, 87, 132}, "m/44'/0'/0'/0/0", false},
		"02877d5614bbb8f602e914d4f4bc0a244479c8790c72f3f09335a53ec53f1bdb3b",
	},
}

var btcAdapterAddressTests = []btcAdapterAddressPair{
	{
		btcAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/0'/0'/0/0", false},
		"13aCiCD1vifQPZRW6kM7iGGFCgzLg2d7iE",
	}, {
		btcAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/1'/0'/0/0", true},
		"n1k5VnD5w71L4WEYw57xBx2VyVZT6hycRD",
	}, {
		btcAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/0'/0'/0/0", false},
		"1DesKpbEh5DfmwsKu5wEAYKFfAoK1Twm4E",
	}, {
		btcAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/0'/0'/0/1", false},
		"1CGmcRnRon1mS9EBYu3p6BLSfR9JLsLDb1",
	}, {
		btcAdapterCommonInput{[]byte{225, 124, 247, 221, 250, 142, 182, 27, 152, 9, 191, 103, 209, 22, 144, 103, 190, 21, 157, 44, 240, 157, 1, 164, 162, 19, 87, 86, 38, 60, 13, 87, 184, 37, 27, 240, 171, 41, 33, 2, 250, 173, 64, 47, 93, 228, 13, 240, 96, 4, 156, 74, 21, 5, 184, 71, 163, 191, 74, 51, 48, 113, 87, 132}, "m/44'/0'/0'/0/0", false},
		"1W8oTwznyezBuV8MUfCAtLQwnUowf1c2J",
	},
}

var btcAdapterSignatureTests = []btcAdapterSignaturePair{
	{
		btcAdapterCommonInput{[]byte{231, 123, 253, 32, 91, 177, 57, 41, 83, 210, 141, 254, 70, 155, 155, 209, 146, 239, 121, 33, 115, 236, 10, 103, 15, 213, 59, 14, 171, 113, 96, 133, 224, 169, 71, 197, 252, 254, 148, 145, 81, 131, 232, 109, 136, 38, 30, 96, 164, 19, 58, 42, 207, 81, 15, 139, 154, 151, 104, 139, 171, 132, 186, 122}, "m/44'/0'/0'/0/0", false},
		`{"inputs":[{"txhash":"81b4c832d70cb56ff957589752eb412a4cab78a25a8fc52d6a09e5bd4404d48a","vout":0}],"outputs":[{"address":"1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa","amount":91234}]}`,
		"01000000018ad40444bde5096a2dc58f5aa278ab4c2a41eb52975857f96fb50cd732c8b481000000006a473044022064deb4f6bd3d283368e0eba6ac00f19a7412d01c8c9ff729bd30630bb0c0592502200134f0badc1796dcc7df13932f18a66454ac8558cafe3cfc82ca6f0b0200fd9b012103cd11c3e23a78a041c004ca3410575b688147ddecdf3e5931e0dda23192c8dcc7ffffffff0162640100000000001976a914c8e90996c7c6080ee06284600c684ed904d14c5c88ac00000000",
	},
}

func TestBTCPrivateKey(t *testing.T) {
	for _, pair := range btcAdapterPrivateKeyTests {
		adapter := adapter.NewBitcoinAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

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

func TestBTCPublicKey(t *testing.T) {
	for _, pair := range btcAdapterPublicKeyTests {
		adapter := adapter.NewBitcoinAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

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

func TestBTCAddress(t *testing.T) {
	for _, pair := range btcAdapterAddressTests {
		adapter := adapter.NewBitcoinAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

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

func TestBTCSignature(t *testing.T) {
	for _, pair := range btcAdapterSignatureTests {
		adapter := adapter.NewBitcoinAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

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
