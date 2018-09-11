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

type EthereumAdapter struct {
	baseadapter.BlockchainAdapter
}

func NewEthereumAdapter(seed []byte, derivationPath string) *EthereumAdapter {
	adapter := new(EthereumAdapter)
	adapter.Seed = seed
	adapter.DerivationPath = derivationPath
	adapter.IsDev = false

	return adapter
}

func (e *EthereumAdapter) DerivePrivateKey() (string, error) {
	btcecPrivKey, err := lib.DerivePrivateKey(e.Seed, e.DerivationPath, e.IsDev)
	if err != nil {
		return "", err
	}

	privateKey := crypto.FromECDSA(btcecPrivKey.ToECDSA())
	privateKeyHex := hexutil.Encode(privateKey)[2:]

	e.PrivateKey = privateKeyHex

	return e.PrivateKey, nil
}

func (e *EthereumAdapter) GetBlockchainNetwork() string {
	if e.IsDev {
		return "testnet"
	}
	return "mainnet"
}

func (e *EthereumAdapter) SetEnvironmentToDevelopment() {
	e.IsDev = true
}

func (e *EthereumAdapter) SetEnvironmentToProduction() {
	e.IsDev = false
}

func (e *EthereumAdapter) CreateSignedTransaction(payload lib.IRawTx) (string, error) {
	privateKey, err := crypto.HexToECDSA(e.PrivateKey)
	if err != nil {
		return "", err
	}

	tx, chainID, err := createRawTransaction(payload)
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		return "", err
	}

	ts := types.Transactions{signedTx}
	txHex := fmt.Sprintf("%x", ts.GetRlp(0))

	return txHex, nil
}

func createRawTransaction(p lib.IRawTx) (*types.Transaction, int64, error) {
	data, _ := json.Marshal(p)
	var payload lib.EthereumRawTx
	err := json.Unmarshal(data, &payload)

	if err != nil {
		return nil, 0, err
	}

	// if nonce, value, gaslimnit, gasprice, chainid is negative
	// address is not "" or 0 address
	// TODO: validate data
	if payload.ChainID < 0 || payload.To == "" ||
		!strings.HasPrefix(payload.To, "0x") || len(payload.To) != 42 {
		return nil, 0, errors.New("Invalid payload data")
	}

	return types.NewTransaction(
		payload.Nonce,
		common.HexToAddress(payload.To),
		big.NewInt(int64(payload.Value)),
		payload.GasLimit,
		big.NewInt(int64(payload.GasPrice)),
		[]byte(payload.Data),
	), payload.ChainID, nil
}
