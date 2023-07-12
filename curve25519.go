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
	var y2my1, x2mx1, lambda big.Int

	y2my1.Sub(y2, y1)
	x2mx1.Sub(x2, x1)
	lambda.Div(&y2my1, &x2mx1)
	lambda.Mod(&lambda, p)

	var x3, y3 big.Int

	// x3 = lambda^2 - x1 - x2
	x3.Mul(&lambda, &lambda)
	x3.Sub(&x3, x1)
	x3.Sub(&x3, x2)
	x3.Mod(&x3, p)

	// y3 = lambda * (x1 - x3) - y1
	y3.Sub(x1, &x3)
	y3.Mul(&lambda, &y3)
	y3.Sub(&y3, y1)
	y3.Mod(&y3, p)

	return &x3, &y3
}

func (curve _c25519) ScalarBaseMult(scalar []byte) (*big.Int, *big.Int) {
	u, _ := curve25519.X25519(scalar, curve25519.Basepoint)
	U := toBigInt(u)
	V := calculatev(U)
	return U, V
}

func (curve _c25519) ScalarMult(Bx, _ *big.Int, scalar []byte) (*big.Int, *big.Int) {
	u, _ := curve25519.X25519(scalar, fromBigInt(Bx))
	U := toBigInt(u)
	V := calculatev(U)
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
