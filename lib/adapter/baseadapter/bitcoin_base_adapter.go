package baseadapter

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"gitlab.com/arout/Vault/lib"
)

type BitcoinBaseAdapter struct {
	BlockchainAdapter
}

func (b *BitcoinBaseAdapter) DerivePrivateKey() (string, error) {
	// obatin private key from seed + derivation path
	btcecPrivKey, err := lib.DerivePrivateKey(b.Seed, b.DerivationPath, b.IsDev)
	if err != nil {
		return "", err
	}

	network := &chaincfg.MainNetParams
	if b.IsDev {
		network = &chaincfg.TestNet3Params
	}

	privateWIF, err := btcutil.NewWIF(btcecPrivKey, network, true)
	if err != nil {
		return "", err
	}

	// store private string as internal data
	b.PrivateKey = privateWIF.String()

	return b.PrivateKey, nil
}

func (b *BitcoinBaseAdapter) GetBlockchainNetwork() string {
	if b.IsDev {
		return "testnet"
	}
	return "mainnet"
}

// TODO: check testnet signature
func (b *BitcoinBaseAdapter) SetEnvironmentToDevelopment() {
	b.IsDev = true
}

func (b *BitcoinBaseAdapter) SetEnvironmentToProduction() {
	b.IsDev = false
}

func (b *BitcoinBaseAdapter) CreateSignedTransaction(p lib.IRawTx) (string, error) {
	network := &chaincfg.MainNetParams
	if b.IsDev {
		network = &chaincfg.TestNet3Params
	}

	wif, err := btcutil.DecodeWIF(b.PrivateKey)
	if err != nil {
		return "", err
	}

	transaction := wire.NewMsgTx(wire.TxVersion)
	payload, err := parsePayload(p)
	if err != nil {
		return "", err
	}

	// TODO: add validation for txHash length and address length and
	// index, amount should be >= 0

	// add inputs to transaction
	for _, utxo := range payload.Inputs {
		hash, _ := chainhash.NewHashFromStr(utxo.TxHash)
		out := wire.NewOutPoint(hash, utxo.Vout)
		txIn := wire.NewTxIn(out, nil, nil)

		transaction.AddTxIn(txIn)
	}

	// add outputs to transaction
	for _, out := range payload.Outputs {
		destinationAddress, err := btcutil.DecodeAddress(out.Address, network)
		pkScript, err := txscript.PayToAddrScript(destinationAddress)
		if err != nil {
			return "", err
		}

		txOut := wire.NewTxOut(out.Amount, pkScript)
		transaction.AddTxOut(txOut)
	}

	// sign transaction
	for i := range payload.Inputs {
		sigScript, err := txscript.SignatureScript(transaction, i, nil, txscript.SigHashAll, wif.PrivKey, false)
		if err != nil {
			return "", err
		}

		transaction.TxIn[i].SignatureScript = sigScript
	}

	var signedTx bytes.Buffer
	transaction.Serialize(&signedTx)

	// TODO: add validation as specified in tutorial
	return hex.EncodeToString(signedTx.Bytes()), nil
}

func parsePayload(p lib.IRawTx) (lib.BitcoinRawTx, error) {
	data, _ := json.Marshal(p)
	var payload lib.BitcoinRawTx
	err := json.Unmarshal(data, &payload)

	return payload, err
}
