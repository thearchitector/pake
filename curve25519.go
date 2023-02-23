package pake

import (
	"math/big"
	"sync"

	"golang.org/x/crypto/curve25519"
)

/*

utilities
- RFC 7748 Section 5
- https://www.rfc-editor.org/rfc/rfc7748#section-5
*/

var p, _ = new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819949", 10)
var pow2, pow3, A = big.NewInt(2), big.NewInt(3), big.NewInt(486662)

func swapEndianness(buf []byte) []byte {
	invariant := len(buf) - 1
	for i := 0; i < len(buf)/2; i++ {
		buf[i], buf[invariant-i] = buf[invariant-i], buf[i]
	}
	return buf
}

func toBigInt(point []byte) *big.Int {
	buf := make([]byte, 32)
	copy(buf, point)
	buf[len(buf)-1] &= (1 << 7) - 1
	return new(big.Int).SetBytes(swapEndianness(buf))
}

func fromBigInt(point *big.Int) []byte {
	buf := make([]byte, 32)
	// point.Mod(point, p)
	return swapEndianness(point.FillBytes(buf))
}

/* interface implementations */

type _c25519 struct {
	P *big.Int
}

func (curve _c25519) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	return nil, nil
}

func (curve _c25519) ScalarBaseMult(scalar []byte) (*big.Int, *big.Int) {
	u, _ := curve25519.X25519(scalar, curve25519.Basepoint)
	U := toBigInt(u)
	V := new(big.Int)
	return U, V
}

func (curve _c25519) ScalarMult(Bx, _ *big.Int, scalar []byte) (*big.Int, *big.Int) {
	u, _ := curve25519.X25519(scalar, fromBigInt(Bx))
	U := toBigInt(u)
	V := new(big.Int)
	return U, V
}

func (curve _c25519) IsOnCurve(x, y *big.Int) bool {
	return false
}

/* singleton initialization */

var crv25519 _c25519
var initialize sync.Once

func initCurve25519() {
	crv25519.P = p
}

func Curve25519() EllipticCurve {
	initialize.Do(initCurve25519)
	return crv25519
}
