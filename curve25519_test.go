package pake

import (
	"bytes"
	"math/big"
	"testing"

	"golang.org/x/crypto/curve25519"
)

func TestBigIntConversion(t *testing.T) {
	a := toBigInt(curve25519.Basepoint)
	if a.Cmp(big.NewInt(9)) != 0 {
		t.Error("big integer didn't convert correctly")
	}

	b := fromBigInt(a)
	if bytes.Compare(b, curve25519.Basepoint) != 0 {
		t.Error("big integer transformation failed identity test")
	}
}

func TestScalarMult1(t *testing.T) {
	// input scalar, input u and expected u are from RFC 7748 Section 5.2
	// test vector 1
	// https://www.rfc-editor.org/rfc/rfc7748#section-5.2
	IS, _ := new(big.Int).SetString("31029842492115040904895560451863089656472772604678260265531221036453811406496", 10)
	IU, _ := new(big.Int).SetString("34426434033919594451155107781188821651316167215306631574996226621102155684838", 10)
	// decimal encoding of output u, by func in Section 5
	// https://colab.research.google.com/drive/1HvR2DXaTQhvWFWRQoBUNwgCbHkEmILM5
	expected, _ := new(big.Int).SetString("37325765543539916631701301279660700968428932651319597985674090122993663859395", 10)
	OU, _ := Curve25519().ScalarMult(IU, nil, fromBigInt(IS))
	if expected.Cmp(OU) != 0 {
		t.Error("given output didn't match the expected output #1")
	}

	IS.SetString("35156891815674817266734212754503633747128614016119564763269015315466259359304", 10)
	IU.SetString("8883857351183929894090759386610649319417338800022198945255395922347792736741", 10)
	expected.SetString("39566196721700740701373067725336211924689549479508623342842086701180565506965", 10)
	OU, _ = Curve25519().ScalarMult(IU, nil, fromBigInt(IS))

	if expected.Cmp(OU) != 0 {
		t.Error("given output didn't match the expected output #2")
	}
}

func TestScalarMult2(t *testing.T) {
	// same as above, but test vector 2
	k, _ := new(big.Int).SetString("28948022309329048855892746252171976963317496166410141009864396001978282409992", 10)
	u := big.NewInt(9)
	expected, _ := new(big.Int).SetString("54815864700279561125610391355931320566748822376190344121911385527384361806914", 10)

	var temp big.Int
	for i := 0; i < 1; i++ {
		temp.Set(k)
		k, _ = Curve25519().ScalarMult(u, nil, fromBigInt(k))
		u.Set(&temp)
	}

	if expected.Cmp(k) != 0 {
		t.Error("iterative test failed")
	}
}

func TestIsOnCurve(t *testing.T) {
	// u, v of the curve basepoint
	// these are U(P) and V(P) from RFC 7748 Section 4.1
	// https://www.rfc-editor.org/rfc/rfc7748#section-4.1
	u := big.NewInt(9)
	v, _ := new(big.Int).SetString("14781619447589544791020593568409986887264606134616475288964881837755586237401", 10)
	if !Curve25519().IsOnCurve(u, v) {
		t.Error("the basepoint is not on the curve")
	}

	// 1, 1 isn't on the curve
	u.SetInt64(1)
	v.SetInt64(1)
	if Curve25519().IsOnCurve(u, v) {
		t.Error("the point 1,1 is off the curve but claimed to be on it")
	}
}
