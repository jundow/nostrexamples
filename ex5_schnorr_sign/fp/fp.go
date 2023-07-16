package fp

import (
	"fmt"
	"math/big"
)

var zero *big.Int
var one *big.Int
var four *big.Int

//var pSecp256k *big.Int

func init() {
	fmt.Println("fp init")
	zero = big.NewInt(0)
	one = big.NewInt(1)
	four = big.NewInt(4)
}

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
	tmp := big.NewInt(0)
	FpInv(y, p, tmp)
	FpMul(x, tmp, p, ret)
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

func FpSecp256kSqrt(x *big.Int, p *big.Int, ret *big.Int) (*big.Int, bool) {
	xtmp := big.NewInt(0).Set(x)
	tmp := big.NewInt(0)
	tmp2 := big.NewInt(0)

	tmp.Add(p, one)
	tmp.Div(tmp, four)

	FpPow(xtmp, tmp, p, ret)
	FpMul(ret, ret, p, tmp2)

	//fmt.Println(xtmp, tmp, ret, tmp2)

	if tmp2.Cmp(xtmp) != 0 {
		ret.SetInt64(0)
		return ret, false
	}
	return ret, true
}
