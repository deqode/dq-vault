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
	"gitlab.com/arout/Vault/lib"
	"gitlab.com/arout/Vault/lib/adapter/baseadapter"
)

// EthereumAdapter - Ethereum blockchain transaction adapter
type EthereumAdapter struct {
	baseadapter.BlockchainAdapter
}

// NewEthereumAdapter constructor function for EthereumAdapter
// sets seed, derivation path as internal data
func NewEthereumAdapter(seed []byte, derivationPath string) *EthereumAdapter {
	adapter := new(EthereumAdapter)
	adapter.Seed = seed
	adapter.DerivationPath = derivationPath
	adapter.IsDev = false

	return adapter
}

// DerivePrivateKey Derives derivation path to obtain private key
// checks for errors
func (e *EthereumAdapter) DerivePrivateKey() (string, error) {
	// obatin private key from seed + derivation path
	btcecPrivKey, err := lib.DerivePrivateKey(e.Seed, e.DerivationPath, e.IsDev)
	if err != nil {
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

// CreateSignedTransaction creates and signs raw transaction from payload data + private key
func (e *EthereumAdapter) CreateSignedTransaction(payload lib.IRawTx) (string, error) {
	// convert hex to ECDSA private key
	privateKey, err := crypto.HexToECDSA(e.PrivateKey)
	if err != nil {
		return "", err
	}

	// creates raw transaction from payload
	tx, chainID, err := createRawTransaction(payload)
	if err != nil {
		return "", err
	}
	// sign raw transaction using raw transaction + chainId + private key
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		return "", err
	}
	// obtains signed transaction hex
	ts := types.Transactions{signedTx}
	txHex := fmt.Sprintf("%x", ts.GetRlp(0))

	return txHex, nil
}

// generates raw transaction from payload
// returns raw transaction + chainId + error (if any)
func createRawTransaction(p lib.IRawTx) (*types.Transaction, int64, error) {
	data, _ := json.Marshal(p)
	var payload lib.EthereumRawTx
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return nil, 0, err
	}

	// TODO: add validations
	// validate payload data
	if payload.ChainID < 0 || payload.To == "" ||
		!strings.HasPrefix(payload.To, "0x") || len(payload.To) != 42 {
		return nil, 0, errors.New("Invalid payload data")
	}
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
