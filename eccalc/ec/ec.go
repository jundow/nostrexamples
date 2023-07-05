package fp

import (
	"math/big"
)

type EC struct {
	P     *big.Int
	A     *big.Int
	B     *big.Int
	Gx    *big.Int
	Gy    *big.Int
	Order *big.Int
}

type ECElement struct {
	ec *EC
	x  *big.Int
	y  *big.Int
}

var Secp256k1 EC
var test EC

func init() {
	Secp256k1.P, _ = big.NewInt(0).SetString("0017", 16)
	Secp256k1.Order, _ = big.NewInt(0).SetString("0000", 16)
	Secp256k1.A, _ = big.NewInt(0).SetString("0001", 16)
	Secp256k1.B, _ = big.NewInt(0).SetString("0001", 16)
	Secp256k1.Gx, _ = new(big.Int).SetString("0000", 16)
	Secp256k1.Gy, _ = new(big.Int).SetString("0001", 16)

	test.P, _ = big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	test.Order, _ = big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	test.A, _ = big.NewInt(0).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
	test.B, _ = big.NewInt(0).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	test.Gx, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	test.Gy, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
}

func (z *ECElement) Add(x1 *ECElement, x2 *ECElement) *ECElement {
	if x1.x.CmpAbs(z.ec.P) == 0 && x2.x.CmpAbs(z.ec.P) == 0 {

	} else if x1.x.CmpAbs(z.ec.P) == 0 {

	} else if x2.x.CmpAbs(z.ec.P) == 0 {

	} else {

	}

}

func (z *ECElement) Double(x1 *ECElement) *ECElement {
	if x1.x.CmpAbs(z.ec.P) == 0 {

	} else {

	}
}

func (z *ECElement) ScalarMul(x1 *ECElement, m *big.Int) *ECElement {
	if x1.x.CmpAbs(z.ec.P) == 0 {

	} else {

	}
}
