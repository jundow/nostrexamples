package fp

import (
	"math/big"
)

func FpAdd(x *big.Int, y *big.Int, p *big.Int, ret *big.Int) *big.Int {
	ret.Add(x, y)
	ret.Mod(ret, p)
	return ret
}

func FpSub(x *big.Int, y *big.Int, p *big.Int, ret *big.Int) *big.Int {
	ret.Sub(x, y)
	ret.Mod(ret, p)
	return ret
}

func FpMul(x *big.Int, y *big.Int, p *big.Int, ret *big.Int) *big.Int {
	ret.Mul(x, y)
	ret.Mod(ret, p)
	return ret
}

func FpInv(x *big.Int, p *big.Int, ret *big.Int) *big.Int {
	zero := big.NewInt(0)

	a := big.NewInt(0).Set(p)
	b := big.NewInt(0).Set(x)
	q := big.NewInt(0)
	r := big.NewInt(0)

	m11 := big.NewInt(1)
	m12 := big.NewInt(0)
	m21 := big.NewInt(0)
	m22 := big.NewInt(1)

	m11n := big.NewInt(1)
	m12n := big.NewInt(0)
	m21n := big.NewInt(0)
	m22n := big.NewInt(1)

	for {
		m11n.Set(m12)
		m12n.Mul(q, m12).Sub(m11, m12n)
		//m12n.Sub(m11, m12n)
		m21n.Set(m22)
		m22n.Mul(q, m22).Sub(m21, m22n)
		//m22n.Sub(m21, m22n)

		q.DivMod(a, b, r)
		a.Set(b)
		b.Set(r)

		//fmt.Println(an, bn, qn, rn, m11n, m12n, m21n, m22n)

		if r.Cmp(zero) == 0 {
			break
		}

		m11.Set(m11n)
		m12.Set(m12n)
		m21.Set(m21n)
		m22.Set(m22n)
	}
	ret.Mod(m12n, p)
	//fmt.Println(m12n, p, ret)
	return ret
}

func FpDiv(x *big.Int, y *big.Int, p *big.Int, ret *big.Int) *big.Int {
	FpInv(y, p, ret)
	FpMul(x, ret, p, ret)
	return ret
}

func FpPow(x *big.Int, y *big.Int, p *big.Int, ret *big.Int) *big.Int {
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
}

func Fpsecp256kSqrt(x *big.INt, ret *big.Int) *big.Int {

}
