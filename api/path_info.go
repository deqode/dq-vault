package api

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"

	"gitlab.com/arout/Vault/lib/rfc6979"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
)

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

// SignBufferSha256 returns Signature of a valid  byte array, Does not validate a transaction digest.
func SignBufferSha256(bufSha256 []byte, privateKey *ecdsa.PrivateKey) []byte {
	var bufSha256Clone = make([]byte, len(bufSha256))
	copy(bufSha256Clone, bufSha256)

	nonce := 0
	// int(time.Now().Unix()
	for {
		// r, s, err := rfc6979.
		r, s, err := rfc6979.SignECDSA(privateKey, bufSha256Clone, sha256.New, nonce)
		nonce++
		if err != nil {
			log.Println(err)
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
		sig := SignBufferSha256(transactionDigest, privKey.ToECDSA())
		sigsHex[index] = hex.EncodeToString(sig)
	}

	return sigsHex
}

// GetDigestHash returns sha256 hash for the input. Expected I/P is of the form {network's chain id + transaction hex}.
func GetDigestHash(digest []byte) []byte {
	hashedDigest := sha256.New()
	hashedDigest.Write(digest)
	return hashedDigest.Sum(nil)
}

func main() {
	digestHex := "3aef3997194701308d57a65214a7a015d98382ab66a9bc0d655de80842b6bfc59ce2d7c019291cad095c0105d1d8070000000000001111050007616e6b69743132010000000001021500e918e7ca8c63e40472c9a2ab28665d06a41e78d034ee1b2ff2b3635d02e5010000010000000001021500e918e7ca8c63e40472c9a2ab28665d06a41e78d034ee1b2ff2b3635d02e5010000021500e918e7ca8c63e40472c9a2ab28665d06a41e78d034ee1b2ff2b3635d02e5050000000000000000"
	wifs := make([]string, 1)
	wifs[0] = "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	digest, err := hex.DecodeString(digestHex)
	if err != nil {
		panic(err)
	}
	hashedDigest := GetDigestHash(digest)
	fmt.Printf("The signature for the specified hex is %v", getSignature(hashedDigest, wifs))
}

// pathInfo corresponds to READ gen/info.
func (b *backend) pathInfo(_ context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	digestHex := "3aef3997194701308d57a65214a7a015d98382ab66a9bc0d655de80842b6bfc59ce2d7c019291cad095c0105d1d8070000000000001111050007616e6b69743132010000000001021500e918e7ca8c63e40472c9a2ab28665d06a41e78d034ee1b2ff2b3635d02e5010000010000000001021500e918e7ca8c63e40472c9a2ab28665d06a41e78d034ee1b2ff2b3635d02e5010000021500e918e7ca8c63e40472c9a2ab28665d06a41e78d034ee1b2ff2b3635d02e5050000000000000000"
	wifs := make([]string, 1)
	wifs[0] = "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	digest, err := hex.DecodeString(digestHex)
	if err != nil {
		panic(err)
	}
	hashedDigest := GetDigestHash(digest)
	sig := getSignature(hashedDigest, wifs)

	return &logical.Response{
		Data: map[string]interface{}{
			"Info": sig,
		},
	}, nil
}

//required signature
// 20132640869edd6fbfc00d569683b7278560c96cc2d9fcba106d50d60187f8621850371c59fdf90712df5b996dc04472198853b13061549d69b5a728dab5bdbae5
// 20132640869edd6fbfc00d569683b7278560c96cc2d9fcba106d50d60187f8621850371c59fdf90712df5b996dc04472198853b13061549d69b5a728dab5bdbae5
