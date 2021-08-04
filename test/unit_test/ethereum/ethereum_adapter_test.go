package tests

import (
	"testing"

	log "github.com/mgutz/logxi/v1"
	"github.com/deqode/dq-vault/lib/adapter"
)

var logger = log.New("tests")

type ethAdapterCommonInput struct {
	seed           []byte
	derivationPath string
	isDev          bool
}

type ethAdapterPrivateKeyPair struct {
	input      ethAdapterCommonInput
	privateKey string
}

type ethAdapterPublicKeyPair struct {
	input     ethAdapterCommonInput
	publicKey string
}

type ethAdapterAddressPair struct {
	input   ethAdapterCommonInput
	address string
}

type ethAdapterSignaturePair struct {
	input     ethAdapterCommonInput
	payload   string
	signature string
}

var ethAdapterPrivateKeyTests = []ethAdapterPrivateKeyPair{
	{
		ethAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/60'/0'/0/0", false},
		"cecc0787b493865ac37897e6c7ea0888b36be06ed3890f34b8e670131e94055a",
	}, {
		ethAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/60'/0'/0/1", false},
		"b42dd93b751379b358c6273846e03900dba113bc9026f20b8c1ccb922562704a",
	}, {
		ethAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/60'/0'/0/0", false},
		"1b918a66eb8dfa8dbe40848df710342dff591f89f05df8a307dd35f981241d81",
	}, {
		ethAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/60'/0'/0/1", false},
		"895877ca76226b981ef38e41e7ef398eb7002f715869d8458d37a07b46980521",
	}, {
		ethAdapterCommonInput{[]byte{225, 124, 247, 221, 250, 142, 182, 27, 152, 9, 191, 103, 209, 22, 144, 103, 190, 21, 157, 44, 240, 157, 1, 164, 162, 19, 87, 86, 38, 60, 13, 87, 184, 37, 27, 240, 171, 41, 33, 2, 250, 173, 64, 47, 93, 228, 13, 240, 96, 4, 156, 74, 21, 5, 184, 71, 163, 191, 74, 51, 48, 113, 87, 132}, "m/44'/60'/0'/0/0", false},
		"754277d517ad57523f877f17a3a870001f44e867db4f7a14057e4394aeeaf909",
	},
}

var ethAdapterPublicKeyTests = []ethAdapterPublicKeyPair{
	{
		ethAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/60'/0'/0/0", false},
		"03291a9e06bfc251f1dc22614e574d90a37b515ff8c85ba11fbada4b8ba89afff1",
	}, {
		ethAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/60'/0'/0/1", false},
		"02eb0c5e31705e890f224f3c897ad59a98788875fc9a25196dab8942ad8ee78ad4",
	}, {
		ethAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/60'/0'/0/0", false},
		"03bb5d735410c4efc3eff34359bdc9370f93a9b3e2f9980f872102a2663343e9db",
	}, {
		ethAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/60'/0'/0/1", false},
		"03e4034db7fb9f779865eb95c108ae0988b8050433490c7e3d900a4710afbb5ba6",
	}, {
		ethAdapterCommonInput{[]byte{225, 124, 247, 221, 250, 142, 182, 27, 152, 9, 191, 103, 209, 22, 144, 103, 190, 21, 157, 44, 240, 157, 1, 164, 162, 19, 87, 86, 38, 60, 13, 87, 184, 37, 27, 240, 171, 41, 33, 2, 250, 173, 64, 47, 93, 228, 13, 240, 96, 4, 156, 74, 21, 5, 184, 71, 163, 191, 74, 51, 48, 113, 87, 132}, "m/44'/60'/0'/0/0", false},
		"0350931a71c3abcd04a4362ba24c667f20382d9f41a90f84ac4afeb309b29a1ea9",
	},
}

var ethAdapterAddressTests = []ethAdapterAddressPair{
	{
		ethAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/60'/0'/0/0", false},
		"0x8Fbc6F544871754c42a0AB007a98355dA2Eea26A",
	}, {
		ethAdapterCommonInput{[]byte{65, 165, 157, 224, 14, 107, 183, 239, 157, 147, 57, 246, 200, 68, 7, 16, 90, 169, 64, 150, 132, 239, 126, 173, 110, 200, 196, 245, 137, 0, 132, 163, 190, 170, 196, 187, 248, 25, 153, 144, 20, 190, 76, 183, 247, 171, 196, 186, 72, 192, 154, 124, 59, 163, 63, 127, 77, 139, 131, 127, 189, 148, 9, 157}, "m/44'/60'/0'/0/1", false},
		"0xea1c19db55c3853387FF67Be1eAb8E1cD5c5a80A",
	}, {
		ethAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/60'/0'/0/0", false},
		"0x9Ef62EB5d86c7a068df75624891CC8D31D3e996a",
	}, {
		ethAdapterCommonInput{[]byte{133, 163, 43, 108, 188, 241, 2, 80, 123, 80, 27, 81, 104, 195, 1, 15, 245, 126, 29, 213, 164, 31, 173, 239, 71, 211, 42, 106, 164, 95, 85, 17, 243, 100, 197, 42, 183, 234, 169, 171, 60, 74, 87, 61, 39, 125, 228, 27, 233, 141, 69, 118, 53, 22, 78, 146, 227, 250, 11, 246, 72, 130, 218, 15}, "m/44'/60'/0'/0/1", false},
		"0x4B83a889595999997a52Fc62790e6aeb92c2F5E3",
	}, {
		ethAdapterCommonInput{[]byte{225, 124, 247, 221, 250, 142, 182, 27, 152, 9, 191, 103, 209, 22, 144, 103, 190, 21, 157, 44, 240, 157, 1, 164, 162, 19, 87, 86, 38, 60, 13, 87, 184, 37, 27, 240, 171, 41, 33, 2, 250, 173, 64, 47, 93, 228, 13, 240, 96, 4, 156, 74, 21, 5, 184, 71, 163, 191, 74, 51, 48, 113, 87, 132}, "m/44'/60'/0'/0/0", false},
		"0xe399678B406803d3698792eD856a9eF2B78EB4E2",
	},
}

var ethAdapterSignatureTests = []ethAdapterSignaturePair{
	{
		ethAdapterCommonInput{[]byte{231, 123, 253, 32, 91, 177, 57, 41, 83, 210, 141, 254, 70, 155, 155, 209, 146, 239, 121, 33, 115, 236, 10, 103, 15, 213, 59, 14, 171, 113, 96, 133, 224, 169, 71, 197, 252, 254, 148, 145, 81, 131, 232, 109, 136, 38, 30, 96, 164, 19, 58, 42, 207, 81, 15, 139, 154, 151, 104, 139, 171, 132, 186, 122}, "m/44'/60'/0'/0/0", false},
		`{"nonce":0,"value":100000000000000000,"gasLimit":21000,"gasPrice":5,"to":"0x95273d64876408E0eDa01a45775Efc2Df6d1CfaC","data":"","chainId":1}`,
		"0xf86780058252089495273d64876408e0eda01a45775efc2df6d1cfac88016345785d8a00008025a046300c200def6f3e96135aa279c80cd2858e02a4ac8e03e93ad01276dc443f7ca00f6515899a2293ca29747c776f874dd2c71c03bcfe6a1801d87cd4ac8e7558af",
	},
}

func TestETHPrivateKey(t *testing.T) {
	for _, pair := range ethAdapterPrivateKeyTests {
		adapter := adapter.NewEthereumAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

		privateKey, _ := adapter.DerivePrivateKey(logger)
		if privateKey != pair.privateKey {
			t.Error(
				"Seed", pair.input.seed,
				"\nPath", pair.input.derivationPath,
				"\nExpected PrivateKey", pair.privateKey,
				"\nGot PrivateKey", privateKey,
			)
		}
	}
}

func TestETHPublicKey(t *testing.T) {
	for _, pair := range ethAdapterPublicKeyTests {
		adapter := adapter.NewEthereumAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

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

func TestETHAddress(t *testing.T) {
	for _, pair := range ethAdapterAddressTests {
		adapter := adapter.NewEthereumAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

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

func TestETHSignature(t *testing.T) {
	for _, pair := range ethAdapterSignatureTests {
		adapter := adapter.NewEthereumAdapter(pair.input.seed, pair.input.derivationPath, pair.input.isDev)

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
