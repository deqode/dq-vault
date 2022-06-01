package adapter

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/deqode/dq-vault/config"
	"github.com/deqode/dq-vault/lib"
	"github.com/deqode/dq-vault/lib/adapter/baseadapter"
	"github.com/deqode/dq-vault/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

// EthereumAdapter - Ethereum blockchain transaction adapter
type EthereumAdapter struct {
	baseadapter.BlockchainAdapter
	zeroAddress string
}

// NewEthereumAdapter constructor function for EthereumAdapter
// sets seed, derivation path as internal data
func NewEthereumAdapter(seed []byte, derivationPath string, isDev bool) *EthereumAdapter {
	adapter := new(EthereumAdapter)
	adapter.Seed = seed
	adapter.DerivationPath = derivationPath
	adapter.IsDev = isDev
	adapter.zeroAddress = "0x0000000000000000000000000000000000000000"

	return adapter
}

// DerivePrivateKey Derives derivation path to obtain private key
// checks for errors
func (e *EthereumAdapter) DerivePrivateKey(backendLogger log.Logger) (string, error) {
	// obatin private key from seed + derivation path
	btcecPrivKey, err := lib.DerivePrivateKey(e.Seed, e.DerivationPath, e.IsDev)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}

	// ECDSA private key to bytes
	privateKey := crypto.FromECDSA(btcecPrivKey.ToECDSA())

	// bytes to hex encoded string
	// excluding "0x" prefix
	privateKeyHex := hexutil.Encode(privateKey)[2:]

	// store private string as internal data
	e.PrivateKey = privateKeyHex

	return e.PrivateKey, nil
}

func (e *EthereumAdapter) DerivePublicKey(logger log.Logger) (string, error) {
	// obatin private key from seed + derivation path
	if _, err := e.DerivePrivateKey(logger); err != nil {
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(e.PrivateKey)
	if err != nil {
		return "", err
	}

	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("Invalid ECDSA public key")
	}

	publicKeyBytes := crypto.CompressPubkey(publicKeyECDSA)
	return hexutil.Encode(publicKeyBytes)[2:], nil
}

func (e *EthereumAdapter) DeriveAddress(logger log.Logger) (string, error) {
	// obatin private key from seed + derivation path
	if _, err := e.DerivePrivateKey(logger); err != nil {
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(e.PrivateKey)
	if err != nil {
		return "", err
	}

	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("Invalid ECDSA public key")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex(), nil
}

// GetBlockchainNetwork returns network config
// default isDev=false i.e. Mainnet
func (e *EthereumAdapter) GetBlockchainNetwork() string {
	if e.IsDev {
		return "testnet"
	}
	return "mainnet"
}

// CreateSignedTransaction creates and signs raw transaction from payload data + private key
func (e *EthereumAdapter) CreateSignedTransaction(payload string, backendLogger log.Logger) (string, error) {
	// convert hex to ECDSA private key
	privateKey, err := crypto.HexToECDSA(e.PrivateKey)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}

	// creates raw transaction from payload
	tx, chainID, err := e.createRawTransaction(payload, backendLogger)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}

	// sign raw transaction using raw transaction + chainId + private key
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}
	// obtains signed transaction hex
	var signedTxBuff bytes.Buffer
	signedTx.EncodeRLP(&signedTxBuff)
	txHex := hexutil.Encode(signedTxBuff.Bytes())

	return txHex, nil
}

// CreateSignature creates and signs hex message from payload data + private key
func (e *EthereumAdapter) CreateSignature(payload string, backendLogger log.Logger) (string, error) {
	// convert hex to ECDSA private key
	privateKey, err := crypto.HexToECDSA(e.PrivateKey)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}

	// TODO: validate payload

	signHash := createSignFromHex(payload)
	signatureBytes, err := crypto.Sign(signHash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// This signature has recovery ID in the last byte instead of v
	// To match the signature of the Go code with JavaScript code, the recovery ID must be replaced by v,
	// v is one byte in size and is just the value in which both signatures differ.
	// The recovery ID has values between 0 and 3, so v has values between 27 and 30 or 0x1b and 0x1e
	// Hence, v = 27 + rid
	// where rid is recovery ID

	signatureBytes[64] += 27
	return hex.EncodeToString(signatureBytes), err

}

// generates raw transaction from payload
// returns raw transaction + chainId + error (if any)
func (e *EthereumAdapter) createRawTransaction(payloadString string, backendLogger log.Logger) (*types.Transaction, *big.Int, error) {

	var payload lib.EthereumRawTx
	if err := json.Unmarshal([]byte(payloadString), &payload); err != nil ||
		reflect.DeepEqual(payload, lib.EthereumRawTx{}) {
		errorMsg := fmt.Sprintf("Unable to decode payload=[%v]", payloadString)

		logger.Log(backendLogger, config.Error, "signature:", errorMsg)
		return nil, nil, errors.New(errorMsg)
	}

	// validate payload data
	valid, txType := validatePayload(payload, e.zeroAddress)
	if !valid {
		logger.Log(backendLogger, config.Error, "signature:", "Invalid payload data")
		return nil, nil, errors.New("Invalid payload data")
	}

	// logging transaction payload info
	logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("type - %v", txType))
	logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("to - %v", payload.To))
	logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("gas limit - %v", payload.GasLimit))
	logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("gas price - %v", payload.GasPrice))
	logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("value - %v", payload.Value))
	logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("data - %v", payload.Data))
	logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("chain id - %v", payload.ChainID))

	// create raw transaction from payload data
	return types.NewTransaction(
		payload.Nonce,
		common.HexToAddress(payload.To),
		payload.Value,
		payload.GasLimit,
		payload.GasPrice,
		common.FromHex(string(payload.Data)),
	), payload.ChainID, nil
}

// validate payload inputs and returns type of
// transaction if payload is valid
func validatePayload(payload lib.EthereumRawTx, zeroAddress string) (bool, string) {
	// Value, chainId, GasPrice should not be negative
	if payload.ChainID.Cmp(big.NewInt(0)) == -1 ||
		payload.Value.Cmp(big.NewInt(0)) == -1 ||
		payload.GasPrice.Cmp(big.NewInt(0)) == -1 {
		return false, ""
	}

	if payload.To == "" && payload.Data != "" {
		return true, "Contract Creation"
	}

	if payload.To != "" {
		if !common.IsHexAddress(payload.To) ||
			!strings.HasPrefix(payload.To, "0x") || len(payload.To) != 42 ||
			payload.To == zeroAddress {
			return false, ""
		}
		transactionType := "Ether Transfer"
		if payload.Data != "" {
			transactionType = "Contract Function Call"
		}

		return true, transactionType
	}
	return false, ""
}

// createSignFromHex generates hashed string to EIP-191: Signed Data Standard Hash
// https://eips.ethereum.org/EIPS/eip-191
func createSignFromHex(hex string) common.Hash {
	hash := common.HexToHash(hex).Bytes()

	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	return crypto.Keccak256Hash([]byte(msg))
}
