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
	P := new(big.Int).SetBytes(swapEndianness(buf))
	P.Mod(P, p)
	return P
}

func fromBigInt(point *big.Int) []byte {
	buf := make([]byte, 32)
	point.Mod(point, p)
	return swapEndianness(point.FillBytes(buf))
}

func calculatev(u *big.Int) *big.Int {
	var rhs, lrhs, mrhs, rrhs big.Int

	lrhs.Exp(u, pow3, nil) // u^3
	mrhs.Mul(u, u)         // u^2
	rrhs.Mul(&mrhs, A)     // u^2*486662

	rhs.Add(&lrhs, &rrhs) // u^3 + u^2*486662
	rhs.Add(&rhs, u)      // u^3 + u^2*486662 + u
	rhs.Mod(&rhs, p)      // u^3 + u^2*486662 + u (mod p)

	// u^3 + u^2*486662 + u (mod p)
	return &rhs
}

/* interface implementations */

type _c25519 struct {
	P *big.Int
}

func (curve _c25519) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	// https://en.wikipedia.org/wiki/Montgomery_curve#Addition
	var y2my1, x2mx1, sqy2my1, sqx2mx1, x3 big.Int

	y2my1.Sub(y2, y1)
	x2mx1.Sub(x2, x1)
	sqy2my1.Mul(&y2my1, &y2my1)
	sqx2mx1.Mul(&x2mx1, &x2mx1)
	x3.Div(&sqy2my1, &sqx2mx1)
	x3.Sub(&x3, A)
	x3.Sub(&x3, x1)
	x3.Sub(&x3, x2)
	x3.Mod(&x3, p)

	var lhs, rhs, y3 big.Int

	lhs.Mul(x1, x1)
	lhs.Add(&lhs, x2)
	lhs.Add(&lhs, A)
	lhs.Mul(&lhs, &y2my1)

	lhs.Div(&lhs, &x2mx1)
	rhs.Mul(&y2my1, &y2my1)
	rhs.Div(&rhs, new(big.Int).Mul(&x2mx1, &x2mx1))
	rhs.Sub(&rhs, y1)
	y3.Sub(&lhs, &rhs)
	y3.Mod(&y3, p)

	return &x3, &y3
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
	lhs := new(big.Int).Exp(y, pow2, p)
	rhs := calculatev(x)
	// v^2 (mod p) = u^3 + u^2*486662 + u (mod p)
	return lhs.Cmp(rhs) == 0
}

/* singleton initialization */

var crv25519 _c25519
var initialize sync.Once

func initCurve25519() {
	crv25519.P = p
}

func Curve25519() _c25519 {
	initialize.Do(initCurve25519)
	return crv25519
}
