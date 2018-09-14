package baseadapter

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	log "github.com/mgutz/logxi/v1"
	"gitlab.com/arout/Vault/config"
	"gitlab.com/arout/Vault/lib"
	"gitlab.com/arout/Vault/logger"
)

type BitcoinBaseAdapter struct {
	BlockchainAdapter
}

func (b *BitcoinBaseAdapter) DerivePrivateKey(logger log.Logger) (string, error) {
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

func (b *BitcoinBaseAdapter) CreateSignedTransaction(p lib.IRawTx, backendLogger log.Logger) (string, error) {
	network := &chaincfg.MainNetParams
	if b.IsDev {
		network = &chaincfg.TestNet3Params
	}

	//decode wif from private key
	wif, err := btcutil.DecodeWIF(b.PrivateKey)
	if err != nil {
		return "", err
	}

	//parse the input payload
	payload, err := parsePayload(p)
	if err != nil {
		return "", err
	}

	if len(payload.Inputs) == 0 || len(payload.Outputs) == 0 {
		logger.Log(backendLogger, config.Error, "signature:", "Invalid payload structure")
		return "", errors.New("Invalid payload structure")
	}

	//generate pubkeyScript from sender's public key
	pubkey := wif.PrivKey.PubKey()
	pubKeyHash := btcutil.Hash160(pubkey.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, network)
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}

	//generating sender's address from sender's public key
	p2pk, _ := btcutil.NewAddressPubKey(pubkey.SerializeCompressed(), network)
	sourceAddress := p2pk.AddressPubKeyHash().EncodeAddress()

	logger.Log(backendLogger, config.Info, "signature:", "from", sourceAddress)

	transaction := wire.NewMsgTx(wire.TxVersion)

	//adding transaction inputs
	for _, utxo := range payload.Inputs {
		//check for valid utxo format
		if len(utxo.Txhash) != 64 {
			logger.Log(backendLogger, config.Error, "signature:", fmt.Sprintf("Invalid UTXO hash - %v", utxo.Txhash))
			return "", fmt.Errorf("Invalid UTXO hash - %v", utxo.Txhash)
		}

		logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("txId %v, vout %v", utxo.Txhash, utxo.Vout))
		hash, _ := chainhash.NewHashFromStr(utxo.Txhash)
		out := wire.NewOutPoint(hash, utxo.Vout)
		txIn := wire.NewTxIn(out, nil, nil)
		transaction.AddTxIn(txIn)
	}

	totalAmount := int64(0)

	//adding transaction outputs
	for _, out := range payload.Outputs {
		if out.Amount < 0 {
			logger.Log(backendLogger, config.Error, "signature:", fmt.Sprintf("Invalid payee amount %v", out.Amount))
			return "", fmt.Errorf("Invalid payee amount %v", out.Amount)
		}

		//check for to payee address validity
		_, _, err := base58.CheckDecode(out.Address)
		if err != nil {
			if err == base58.ErrChecksum {
				logger.Log(backendLogger, config.Error, "signature:", "Payee address checksum mismatch")
				return "", errors.New("Payee address checksum mismatch")
			}
			return "", errors.New("Invalid payee address format")
		}

		if out.Address == sourceAddress {
			//in case of change to be returned, we use our own pkscript
			transaction.AddTxOut(wire.NewTxOut(out.Amount, pkScript))
			totalAmount += out.Amount
		} else {
			destinationAddress, _ := btcutil.DecodeAddress(out.Address, network)
			destinationPkScript, _ := txscript.PayToAddrScript(destinationAddress)
			transaction.AddTxOut(wire.NewTxOut(out.Amount, destinationPkScript))
			totalAmount += out.Amount
		}

		logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("Payee address %v, amount %v", out.Address, out.Amount))
	}

	// Sign the redeeming transaction.
	lookupKey := func(a btcutil.Address) (*btcec.PrivateKey, bool, error) {
		return wif.PrivKey, true, nil
	}
	// Notice that the script database parameter is nil here since it isn't
	// used.  It must be specified when pay-to-script-hash transactions are
	// being signed.
	for i := range payload.Inputs {
		sigScript, err := txscript.SignTxOutput(network, transaction, i, pkScript, txscript.SigHashAll, txscript.KeyClosure(lookupKey), nil, nil)

		if err != nil {
			logger.Log(backendLogger, config.Error, "signature:", err.Error())
			return "", err
		}
		transaction.TxIn[i].SignatureScript = sigScript
	}

	var signedTx bytes.Buffer
	transaction.Serialize(&signedTx)

	// Prove that the transaction has been validly signed by executing the
	// script pair.
	flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
		txscript.ScriptStrictMultiSig |
		txscript.ScriptDiscourageUpgradableNops
	vm, err := txscript.NewEngine(pkScript, transaction, 0, flags, nil, nil, totalAmount)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}
	if err := vm.Execute(); err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}

	return hex.EncodeToString(signedTx.Bytes()), nil
}

func parsePayload(p lib.IRawTx) (lib.BitcoinRawTx, error) {
	data, _ := json.Marshal(p)
	var payload lib.BitcoinRawTx
	err := json.Unmarshal(data, &payload)

	return payload, err
}
