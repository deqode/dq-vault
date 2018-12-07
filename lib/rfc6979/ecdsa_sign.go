package rfc6979

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"hash"
	"math/big"
	//"log"
	//"encoding/hex"
)

var oneInitializer = []byte{0x01}

//HmacSHA256  returns a Hash-based message authentication code
func HmacSHA256(m []byte, k []byte) []byte {
	//return mac(fastsha256.New, crypto.Sha256(m), crypto.Sha256(k))
	mac := hmac.New(sha256.New, k)
	mac.Write(m)
	expectedMAC := mac.Sum(nil)
	return expectedMAC
}

// https://tools.ietf.org/html/rfc6979#section-3.2
func generateSecret(priv *ecdsa.PrivateKey /*q, x *big.Int,*/, alg func() hash.Hash, hash []byte, test func(*big.Int) bool, nonce int) {
	//log.Println("priv=", priv.)
	var hashClone = make([]byte, len(hash))
	copy(hashClone, hash)

	//log.Println("before hashClone=", hex.EncodeToString(hashClone))
	if nonce > 0 {
		//nonce_str := RandStringBytes(nonce)
		nonceA := make([]byte, 4)
		binary.BigEndian.PutUint32(nonceA, uint32(nonce))
		hashClone = append(hashClone, nonceA...)
		//log.Println("(before hash) hashClone=", hex.EncodeToString(hashClone), "nonce_str=", hex.EncodeToString(nonceA))
		hs := sha256.New()
		hs.Write(hashClone)
		hashClone = hs.Sum(nil)
	}
	//log.Println("hashClone=", hex.EncodeToString(hashClone))

	c := priv.PublicKey.Curve
	//N := c.Params().N
	x := priv.D.Bytes()
	q := c.Params().N
	//x := privkey.Bytes()
	//alg := fastsha256.New

	//qlen := q.BitLen()
	//holen := alg().Size()

	//rolen := (qlen + 7) >> 3
	//bx := append(int2octets(x, rolen), bits2octets(hash, curve, rolen)...)

	//log.Println("bx=", hex.EncodeToString(bx))

	// Step B
	v := bytes.Repeat(oneInitializer, 32)

	// Step C (Go zeroes the all allocated memory)
	k := make([]byte, 32)

	// Step D
	//k = mac(alg, k, append(append(append(v, 0x00), bx...), hash... ))
	m := append(append(append(v, 0x00), x...), hashClone...)
	//log.Println("m", hex.EncodeToString(m))
	//log.Println("k", hex.EncodeToString(k))
	k = HmacSHA256(m, k)
	//log.Println("Step D", hex.EncodeToString(k))

	// Step E
	//v = mac(alg, k, v)
	v = HmacSHA256(v, k)
	//log.Println("Step E", hex.EncodeToString(v))

	// Step F
	//k = mac(alg, k, append(append(append(v, 0x01), bx...), hash...))
	k = HmacSHA256(append(append(append(v, 0x01), x...), hashClone...), k)
	//log.Println("Step F", hex.EncodeToString(k))

	// Step G
	//v = mac(alg, k, v)
	v = HmacSHA256(v, k)
	//log.Println("Step G", hex.EncodeToString(v))

	// Step H1/H2a, ignored as tlen === qlen (256 bit)
	// Step H2b
	v = HmacSHA256(v, k)
	//log.Println("Step H2b", hex.EncodeToString(v))

	//if (nonce.Cmp(big.NewInt(0)) != 0) {
	//	alg := sha256.New()
	//	alg.Write(hash)
	//	alg.Write(nonce.Bytes())
	//	hash = alg.Sum(nil)
	//}
	//
	//qlen := q.BitLen()
	//holen := alg().Size()
	//rolen := (qlen + 7) >> 3
	//bx := append(int2octets(x, rolen), bits2octets(hash, q, qlen, rolen)...)
	//
	//// Step B
	//v := bytes.Repeat([]byte{0x01}, holen)
	//
	//// Step C
	//k := bytes.Repeat([]byte{0x00}, holen)
	//
	//// Step D
	//k = mac(alg, k, append(append(v, 0x00), bx...), k)
	//
	//// Step E
	//v = mac(alg, k, v, v)
	//
	//// Step F
	//k = mac(alg, k, append(append(v, 0x01), bx...), k)
	//
	//// Step G
	//v = mac(alg, k, v, v)

	//////////////////////

	var T = hashToInt(v, c)
	//log.Println("T", hex.EncodeToString(T.Bytes()))

	// Step H3, repeat until T is within the interval [1, n - 1]
	for T.Sign() <= 0 || T.Cmp(q) >= 0 || !test(T) {

		//k = crypto.HmacSHA256(Buffer.concat([v, new Buffer([0])]), k);
		k = HmacSHA256(append(v, 0x00), k)

		//v = crypto.HmacSHA256(v, k);
		v = HmacSHA256(v, k)

		// Step H1/H2a, again, ignored as tlen === qlen (256 bit)
		// Step H2b again
		//v = crypto.HmacSHA256(v, k);
		v = HmacSHA256(v, k)

		//T = BigInteger.fromBuffer(v);
		T = hashToInt(v, c)

		//log.Println("T", hex.EncodeToString(T.Bytes()))
	}

	//return T;

	//// Step H
	//for {
	//	// Step H1
	//	var t []byte
	//
	//	// Step H2
	//	for len(t) < qlen/8 {
	//		v = mac(alg, k, v, v)
	//		t = append(t, v...)
	//	}
	//
	//	// Step H3
	//	secret := bits2int(t, qlen)
	//	log.Println("secret", hex.EncodeToString(secret.Bytes()))
	//	if secret.Cmp(one) >= 0 && secret.Cmp(q) < 0 && test(secret) {
	//		return
	//	}
	//	k = mac(alg, k, append(v, 0x00), k)
	//	v = mac(alg, k, v, v)
	//}
}

// SignECDSA signs an arbitrary length hash (which should be the result of
// hashing a larger message) using the private key, priv. It returns the
// signature as a pair of integers.
//
// Note that FIPS 186-3 section 4.6 specifies that the hash should be truncated
// to the byte-length of the subgroup. This function does not perform that
// truncation itself.
func SignECDSA(priv *ecdsa.PrivateKey, hash []byte, alg func() hash.Hash, nonce int) (r, s *big.Int, err error) {
	c := priv.PublicKey.Curve
	N := c.Params().N

	//log.Println("e=", hex.EncodeToString(hash)) 		 // OK
	//log.Println("N=", hex.EncodeToString(N.Bytes()))	 // OK

	var hashClone = make([]byte, len(hash))
	copy(hashClone, hash)

	//log.Println("generateSecret- nonce=", nonce)
	generateSecret(priv /* N, priv.D, */, alg, hashClone, func(k *big.Int) bool {
		inv := new(big.Int).ModInverse(k, N)
		r, _ = priv.Curve.ScalarBaseMult(k.Bytes())
		r.Mod(r, N)

		if r.Sign() == 0 {
			//log.Println("r.Sign() == 0")
			return false
		}

		e := hashToInt(hashClone, c)
		s = new(big.Int).Mul(priv.D, r)
		s.Add(s, e)
		s.Mul(s, inv)
		s.Mod(s, N)

		if s.Sign() == 0 {
			//log.Println("s.Sign() == 0")
			return false
		}

		return true
	}, nonce)

	//log.Println("enforce low S values, see bip62: 'low s values in signatures'");
	// enforce low S values, see bip62: 'low s values in signatures'
	NOverTwo := new(big.Int).Div(N, big.NewInt(2))
	if s.Cmp(NOverTwo) > 0 {
		s = new(big.Int).Sub(N, s)
	}

	return
}

// copied from crypto/ecdsa
func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	var hashClone = make([]byte, len(hash))
	copy(hashClone, hash)

	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hashClone) > orderBytes {
		hashClone = hashClone[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hashClone)
	excess := len(hashClone)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}
