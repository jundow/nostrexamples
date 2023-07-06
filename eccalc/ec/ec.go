package fp

import (
	"eccalc/fp"
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

func NewECElement(ec *EC, gx, gy *big.Int) *ECElement {

}

func (gr *ECElement) Add(g1, g2 *ECElement) *ECElement {
	if g1.x.CmpAbs(gr.ec.P) == 0 && g2.x.CmpAbs(gr.ec.P) == 0 {
		gr.x.Set(gr.ec.P)
		gr.y.Set(gr.ec.P)
	} else if g1.x.CmpAbs(gr.ec.P) == 0 {
		gr.x.Set(g2.x)
		gr.y.Set(g2.y)
	} else if g2.x.CmpAbs(gr.ec.P) == 0 {
		gr.x.Set(g1.x)
		gr.y.Set(g1.y)
	} else {
		A := big.NewInt(0)
		B := big.NewInt(0)
		C := big.NewInt(0)
		ret := big.NewInt(0)

		lambda := big.NewInt(0)

		fp.FpSub(g2.y, g1.y, gr.ec.P, A)
		fp.FpSub(g2.x, g1.x, gr.ec.P, B)
		fp.FpDiv(A, B, gr.ec.P, lambda)

		fp.FpMul(lambda, lambda, gr.ec.P, ret)
		fp.FpSub(ret, g1.x, gr.ec.P, ret)
		fp.FpSub(ret, g2.x, gr.ec.P, gr.x)

		fp.FpSub(g1.x, gr.x, gr.ec.P, C)
		fp.FpMul(lambda, C, gr.ec.P, ret)
		fp.FpSub(ret, g1.y, gr.ec.P, gr.y)
	}
	return gr
}

func (gr *ECElement) Double(g1 *ECElement) *ECElement {
	if g1.x.CmpAbs(gr.ec.P) == 0 {
		gr.x.Set(gr.ec.P)
		gr.y.Set(gr.ec.P)
	} else {
		two := big.NewInt(2)

		C := big.NewInt(0)
		lambda := big.NewInt(3)

		fp.FpMul(lambda, g1.x, gr.ec.P, lambda)
		fp.FpMul(lambda, g1.x, gr.ec.P, lambda)
		fp.FpAdd(lambda, gr.ec.A, gr.ec.P, lambda)
		fp.FpDiv(lambda, g1.y, gr.ec.P, lambda)
		fp.FpDiv(lambda, two, gr.ec.P, lambda)

		ret := big.NewInt(0)

		ret.Set(lambda)
		fp.FpMul(ret, lambda, gr.ec.P, ret)
		fp.FpSub(ret, g1.x, gr.ec.P, ret)
		fp.FpSub(ret, g1.x, gr.ec.P, gr.x)

		fp.FpSub(g1.x, gr.x, gr.ec.P, C)
		fp.FpMul(lambda, C, gr.ec.P, ret)
		fp.FpSub(ret, g1.y, gr.ec.P, gr.y)
	}
	return gr
}

func (gr *ECElement) ScalarMul(g1 *ECElement, m *big.Int) *ECElement {
	if g1.x.CmpAbs(gr.ec.P) == 0 {

	} else {
		/*
			xtmp := big.NewInt(0).Set(x)
			ret.SetInt64(1)
			//fmt.Println(y.BitLen())
			for i := 0; i < y.BitLen(); i++ {
				if y.Bit(i) > 0 {
					ret.Mul(ret, xtmp)
					ret.Mod(ret, p)
				}
				//fmt.Println("Debug", i, y.Bit(i), xtmp, ret)
				xtmp.Mul(xtmp, xtmp)
				xtmp.Mod(xtmp, p)
			}
			return ret
		*/

	}
}
