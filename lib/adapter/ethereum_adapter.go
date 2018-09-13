package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/mgutz/logxi/v1"
	"gitlab.com/arout/Vault/lib"
	"gitlab.com/arout/Vault/lib/adapter/baseadapter"
)

// EthereumAdapter - Ethereum blockchain transaction adapter
type EthereumAdapter struct {
	baseadapter.BlockchainAdapter
	zeroAddress string
}

// NewEthereumAdapter constructor function for EthereumAdapter
// sets seed, derivation path as internal data
func NewEthereumAdapter(seed []byte, derivationPath string) *EthereumAdapter {
	adapter := new(EthereumAdapter)
	adapter.Seed = seed
	adapter.DerivationPath = derivationPath
	adapter.IsDev = false
	adapter.zeroAddress = "0x0000000000000000000000000000000000000000"

	return adapter
}

// DerivePrivateKey Derives derivation path to obtain private key
// checks for errors
func (e *EthereumAdapter) DerivePrivateKey(logger log.Logger) (string, error) {
	// obatin private key from seed + derivation path
	btcecPrivKey, err := lib.DerivePrivateKey(e.Seed, e.DerivationPath, e.IsDev)
	if err != nil {
		logger.Info(fmt.Sprintf("\n[ERROR ] signature: %v", err))
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

// GetBlockchainNetwork returns network config
// default isDev=false i.e. Mainnet
func (e *EthereumAdapter) GetBlockchainNetwork() string {
	if e.IsDev {
		return "testnet"
	}
	return "mainnet"
}

// SetEnvironmentToDevelopment sets environment to Development
func (e *EthereumAdapter) SetEnvironmentToDevelopment() {
	e.IsDev = true
}

// SetEnvironmentToProduction sets environment to Mainnet
func (e *EthereumAdapter) SetEnvironmentToProduction() {
	e.IsDev = false
}

// TODO: verify in Dev mode

// CreateSignedTransaction creates and signs raw transaction from payload data + private key
func (e *EthereumAdapter) CreateSignedTransaction(payload lib.IRawTx, logger log.Logger) (string, error) {
	// convert hex to ECDSA private key
	privateKey, err := crypto.HexToECDSA(e.PrivateKey)
	if err != nil {
		logger.Info(fmt.Sprintf("\n[ERROR ] signature: %v", err))
		return "", err
	}

	// creates raw transaction from payload
	tx, chainID, err := e.createRawTransaction(payload, logger)
	if err != nil {
		logger.Info(fmt.Sprintf("\n[ERROR ] signature: %v", err))
		return "", err
	}
	// sign raw transaction using raw transaction + chainId + private key
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		logger.Info(fmt.Sprintf("\n[ERROR ] signature: %v", err))
		return "", err
	}
	// obtains signed transaction hex
	ts := types.Transactions{signedTx}
	txHex := hexutil.Encode(ts.GetRlp(0))

	return txHex, nil
}

// generates raw transaction from payload
// returns raw transaction + chainId + error (if any)
func (e *EthereumAdapter) createRawTransaction(p lib.IRawTx, logger log.Logger) (*types.Transaction, int64, error) {
	data, _ := json.Marshal(p)
	var payload lib.EthereumRawTx
	err := json.Unmarshal(data, &payload)
	if err != nil {
		logger.Info(fmt.Sprintf("\n[ERROR ] signature: %v", err))
		return nil, 0, err
	}

	// validate payload data
	valid, txType := validatePayload(payload, e.zeroAddress)
	if !valid {
		logger.Info(fmt.Sprintf("\n[ERROR ] signature: Invalid payload data"))
		return nil, 0, errors.New("Invalid payload data")
	}

	// logging transaction payload info
	logger.Info(fmt.Sprintf("\n[INFO ] signature: type - %v", txType))
	logger.Info(fmt.Sprintf("\n[INFO ] signature: to - %v", payload.To))
	logger.Info(fmt.Sprintf("\n[INFO ] signature: gas limit - %v", payload.GasLimit))
	logger.Info(fmt.Sprintf("\n[INFO ] signature: gas price - %v", payload.GasPrice))
	logger.Info(fmt.Sprintf("\n[INFO ] signature: value - %v", payload.Value))
	logger.Info(fmt.Sprintf("\n[INFO ] signature: data - %v", payload.Data))
	logger.Info(fmt.Sprintf("\n[INFO ] signature: chain id - %v", payload.ChainID))

	// create raw transaction from payload data
	return types.NewTransaction(
		payload.Nonce,
		common.HexToAddress(payload.To),
		big.NewInt(int64(payload.Value)),
		payload.GasLimit,
		big.NewInt(int64(payload.GasPrice)),
		[]byte(payload.Data),
	), payload.ChainID, nil
}

// validate payload inputs and returns type of
// transaction if payload is valid
// TODO: improve this
func validatePayload(payload lib.EthereumRawTx, zeroAddress string) (bool, string) {
	if payload.ChainID < 0 {
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
