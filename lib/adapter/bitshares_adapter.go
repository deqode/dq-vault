package adapter

// WIP bitshares adapter

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/mgutz/logxi/v1"
	"gitlab.com/arout/Vault/config"
	"gitlab.com/arout/Vault/lib"
	"gitlab.com/arout/Vault/lib/adapter/baseadapter"
	"gitlab.com/arout/Vault/lib/rfc6979"
	"gitlab.com/arout/Vault/logger"
	"golang.org/x/crypto/ripemd160"
)

// BitsharesAdapter - Ethereum blockchain transaction adapter
type BitsharesAdapter struct {
	baseadapter.BlockchainAdapter
	zeroAddress string
}

// NewBitsharesAdapter constructor function for BitsharesAdapter
// sets seed, derivation path as internal data
func NewBitsharesAdapter(seed []byte, derivationPath string, isDev bool) *BitsharesAdapter {
	adapter := new(BitsharesAdapter)
	adapter.Seed = seed
	adapter.DerivationPath = config.BitsharesDerivationPath
	adapter.IsDev = isDev
	adapter.zeroAddress = "0x0000000000000000000000000000000000000000"

	return adapter
}

// DerivePrivateKey Derives derivation path to obtain private key
// checks for errors
func (e *BitsharesAdapter) DerivePrivateKey(backendLogger log.Logger) (string, error) {
	// obatin private key from seed + derivation path
	btcecPrivKey, err := lib.DerivePrivateKey(e.Seed, e.DerivationPath, e.IsDev)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}

	network := &chaincfg.MainNetParams

	privateWIF, err := btcutil.NewWIF(btcecPrivKey, network, false)
	if err != nil {
		return "", err
	}

	// store private string as internal data
	e.PrivateKey = privateWIF.String()
	return e.PrivateKey, nil
}

// DerivePublicKey returns the public key for BTS format.
func (e *BitsharesAdapter) DerivePublicKey(logger log.Logger) (string, error) {
	// obatin private key from seed + derivation path
	if _, err := e.DerivePrivateKey(logger); err != nil {
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(e.PrivateKey)
	if err != nil {
		return "", err
	}

	mdHash := ripemd160.New()
	mdHash.Write(toBytes(privateKey.PublicKey))
	checkSum := mdHash.Sum(nil)
	appendedCS := append(toBytes(privateKey.PublicKey), checkSum[0:4]...)
	publicKey := "BTS" + base58.Encode(appendedCS)

	return publicKey, nil
}

// DeriveAddress Address in Bitsahres can be avoided for most cases by using account names instead.
// TODO : ADD a proper address generation logic.
func (e *BitsharesAdapter) DeriveAddress(logger log.Logger) (string, error) {
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
func (e *BitsharesAdapter) GetBlockchainNetwork() string {
	if e.IsDev {
		return "testnet"
	}
	return "mainnet"
}

// CreateSignedTransaction creates and signs raw transaction from transaction digest + private key
// TODO :- INCLUDE BTS TRANSACTION DIGEST SIGNATURE LOGIC.
func (e *BitsharesAdapter) CreateSignedTransaction(payload lib.IRawTx, backendLogger log.Logger) (string, error) {
	if _, err := e.DerivePrivateKey(backendLogger); err != nil {
		return "", err
	}

	wifs := make([]string, 1)
	wifs[0] = e.PrivateKey

	// creates raw transaction from payload
	digestString, err := e.createRawTransaction(payload, backendLogger)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}

	digest, err := hex.DecodeString(digestString)
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}
	hashedDigest := getDigestHash(digest)
	signature := getSignature(hashedDigest, wifs)

	return signature[0], err
}

// generates raw transaction from payload
// returns raw transaction + chainId + error (if any)
func (e *BitsharesAdapter) createRawTransaction(p lib.IRawTx, backendLogger log.Logger) (string, error) {
	data, _ := json.Marshal(p)
	var payload lib.BitsharesRawTx
	err := json.Unmarshal(data, &payload) // payload is now a BitsharesRawTx
	if err != nil {
		logger.Log(backendLogger, config.Error, "signature:", err.Error())
		return "", err
	}

	// validate payload data
	valid := validateBTSPayload(payload)
	if !valid {
		logger.Log(backendLogger, config.Error, "signature:", "Invalid payload data")
		return "", errors.New("Invalid payload data")
	}

	// logging transaction payload info
	logger.Log(backendLogger, config.Info, "signature:", fmt.Sprintf("to - %v", payload.TransactionDigest))

	// create raw transaction from payload data
	return payload.TransactionDigest, nil
}

// validate payload inputs and returns type of
// transaction if payload is valid
func validateBTSPayload(payload lib.BitsharesRawTx) bool {
	// _, err := strconv.ParseUint(payload.TransactionDigest, 16, 64)
	// if err != nil {
	// 	// a valid hex string is not sent in the payload.
	// 	return false
	// }
	return true
}

func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}

func decompressPoint(curve *secp256k1.KoblitzCurve, x *big.Int, ybit bool) (*big.Int, error) {
	// TODO: This will probably only work for secp256k1 due to
	// optimizations.

	// Y = +-sqrt(x^3 + B)
	x3 := new(big.Int).Mul(x, x)
	x3.Mul(x3, x)
	x3.Add(x3, curve.Params().B)

	// now calculate sqrt mod p of x2 + B
	// This code used to do a full sqrt based on tonelli/shanks,
	// but this was replaced by the algorithms referenced in
	// https://bitcointalk.org/index.php?topic=162805.msg1712294#msg1712294
	y := new(big.Int).Exp(x3, curve.QPlus1Div4(), curve.Params().P)

	if ybit != isOdd(y) {
		y.Sub(curve.Params().P, y)
	}
	if ybit != isOdd(y) {
		return nil, fmt.Errorf("ybit doesn't match oddness")
	}
	return y, nil
}

func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}

func recoverKeyFromSignature(curve *secp256k1.KoblitzCurve, sig *secp256k1.Signature, msg []byte, iter int, doChecks bool) (*secp256k1.PublicKey, error) {
	// 1.1 x = (n * i) + r
	Rx := new(big.Int).Mul(curve.Params().N,
		new(big.Int).SetInt64(int64(iter/2)))
	Rx.Add(Rx, sig.R)
	if Rx.Cmp(curve.Params().P) != -1 {
		return nil, errors.New("calculated Rx is larger than curve P")
	}

	// convert 02<Rx> to point R. (step 1.2 and 1.3). If we are on an odd
	// iteration then 1.6 will be done with -R, so we calculate the other
	// term when uncompressing the point.
	Ry, err := decompressPoint(curve, Rx, iter%2 == 1)
	if err != nil {
		return nil, err
	}

	// 1.4 Check n*R is point at infinity
	if doChecks {
		nRx, nRy := curve.ScalarMult(Rx, Ry, curve.Params().N.Bytes())
		if nRx.Sign() != 0 || nRy.Sign() != 0 {
			return nil, errors.New("n*R does not equal the point at infinity")
		}
	}

	// 1.5 calculate e from message using the same algorithm as ecdsa
	// signature calculation.
	e := hashToInt(msg, curve)

	// Step 1.6.1:
	// We calculate the two terms sR and eG separately multiplied by the
	// inverse of r (from the signature). We then add them to calculate
	// Q = r^-1(sR-eG)
	invr := new(big.Int).ModInverse(sig.R, curve.Params().N)

	// first term.
	invrS := new(big.Int).Mul(invr, sig.S)
	invrS.Mod(invrS, curve.Params().N)
	sRx, sRy := curve.ScalarMult(Rx, Ry, invrS.Bytes())

	// second term.
	e.Neg(e)
	e.Mod(e, curve.Params().N)
	e.Mul(e, invr)
	e.Mod(e, curve.Params().N)
	minuseGx, minuseGy := curve.ScalarBaseMult(e.Bytes())

	// TODO: this would be faster if we did a mult and add in one
	// step to prevent the jacobian conversion back and forth.
	Qx, Qy := curve.Add(sRx, sRy, minuseGx, minuseGy)

	return &secp256k1.PublicKey{
		Curve: curve,
		X:     Qx,
		Y:     Qy,
	}, nil
}

func toBytes(pub ecdsa.PublicKey) []byte {
	x := pub.X.Bytes()

	/* Pad X to 32-bytes */
	paddedX := append(bytes.Repeat([]byte{0x00}, 32-len(x)), x...)

	/* Add prefix 0x02 or 0x03 depending on ylsb */
	if pub.Y.Bit(0) == 0 {
		return append([]byte{0x02}, paddedX...)
	}

	return append([]byte{0x03}, paddedX...)
}

// SignBufferSha256 returns Signature of a valid  byte array, Does not validate a transaction digest.
func signBufferSha256(bufSha256 []byte, privateKey *ecdsa.PrivateKey) []byte {
	var bufSha256Clone = make([]byte, len(bufSha256))
	copy(bufSha256Clone, bufSha256)

	nonce := 0
	// int(time.Now().Unix()
	for {
		// r, s, err := rfc6979.
		r, s, err := rfc6979.SignECDSA(privateKey, bufSha256Clone, sha256.New, nonce)
		nonce++
		if err != nil {
			// log.Println(err)
			return nil
		}

		ecsignature := &secp256k1.Signature{R: r, S: s}

		der := ecsignature.Serialize()
		lenR := der[3]
		lenS := der[5+lenR]

		if lenR == 32 && lenS == 32 {
			// bitcoind checks the bit length of R and S here. The ecdsa signature
			// algorithm returns R and S mod N therefore they will be the bitsize of
			// the curve, and thus correctly sized.
			key := (*secp256k1.PrivateKey)(privateKey)
			curve := secp256k1.S256()
			maxCounter := 4 //maxCounter := (curve.H+1)*2
			for i := 0; i < maxCounter; i++ {
				//for i := 0; i < (curve.H+1)*2; i++ {
				//for i := 0; ;i++ {
				pk, err := recoverKeyFromSignature(curve, ecsignature, bufSha256Clone, i, true)

				if err == nil && pk.X.Cmp(key.X) == 0 && pk.Y.Cmp(key.Y) == 0 {
					//result := make([]byte, 1, 2*curve.byteSize+1)
					byteSize := curve.BitSize / 8
					result := make([]byte, 1, 2*byteSize+1)
					result[0] = 27 + byte(i)
					if true { // isCompressedKey
						result[0] += 4
					}
					// Not sure this needs rounding but safer to do so.
					curvelen := (curve.BitSize + 7) / 8

					// Pad R and S to curvelen if needed.
					bytelen := (ecsignature.R.BitLen() + 7) / 8
					if bytelen < curvelen {
						result = append(result, make([]byte, curvelen-bytelen)...)
					}
					result = append(result, ecsignature.R.Bytes()...)

					bytelen = (ecsignature.S.BitLen() + 7) / 8
					if bytelen < curvelen {
						result = append(result, make([]byte, curvelen-bytelen)...)
					}
					result = append(result, ecsignature.S.Bytes()...)

					return result
				}

			}

		}
	}
}

func getSignature(transactionDigest []byte, wifs []string) []string {
	privKeys := make([]*secp256k1.PrivateKey, len(wifs))
	for index, wif := range wifs {
		w, err := btcutil.DecodeWIF(wif) // check it
		if err != nil {
			panic(err)
		}
		privKeys[index] = w.PrivKey
	}

	sigsHex := make([]string, len(privKeys))
	for index, privKey := range privKeys {
		sig := signBufferSha256(transactionDigest, privKey.ToECDSA())
		sigsHex[index] = hex.EncodeToString(sig)
	}

	return sigsHex
}

// getDigestHash returns sha256 hash for the input. Expected I/P is of the form {network's chain id + transaction hex}.
func getDigestHash(digest []byte) []byte {
	hashedDigest := sha256.New()
	hashedDigest.Write(digest)
	return hashedDigest.Sum(nil)
}
