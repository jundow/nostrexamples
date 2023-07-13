package ec

import (
	"eccalc/fp"
	"errors"
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
	Ec *EC
	X  *big.Int
	Y  *big.Int
}

var Secp256k1 EC
var Test EC
var zero *big.Int
var two *big.Int

func init() {
	Test.P, _ = big.NewInt(0).SetString("0017", 16)
	Test.Order, _ = big.NewInt(0).SetString("0000", 16)
	Test.A, _ = big.NewInt(0).SetString("0001", 16)
	Test.B, _ = big.NewInt(0).SetString("0001", 16)
	Test.Gx, _ = new(big.Int).SetString("0000", 16)
	Test.Gy, _ = new(big.Int).SetString("0001", 16)

	Secp256k1.P, _ = big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	Secp256k1.Order, _ = big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	Secp256k1.A, _ = big.NewInt(0).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
	Secp256k1.B, _ = big.NewInt(0).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	Secp256k1.Gx, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	Secp256k1.Gy, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)

	zero = big.NewInt(0)
	two = big.NewInt(2)
}

func NewECElement(ec *EC, gx, gy *big.Int) *ECElement {
	newecelm := new(ECElement)
	newecelm.Ec = ec
	newecelm.X = big.NewInt(0).Set(gx)
	newecelm.Y = big.NewInt(0).Set(gy)
	return newecelm
}

func (g *ECElement) IsUnitElement() bool {
	return g.Y.CmpAbs(g.Ec.P) == 0
}

func (gr *ECElement) SetUnitElement() *ECElement {
	gr.X.Set(gr.Ec.P)
	gr.Y.Set(gr.Ec.P)
	return gr
}

func (gr *ECElement) Set(g *ECElement) *ECElement {
	gr.X.Set(g.X)
	gr.Y.Set(g.Y)
	return gr
}

func (gr *ECElement) Add(g1, g2 *ECElement) *ECElement {
	if g1.IsUnitElement() && g2.IsUnitElement() {
		gr.SetUnitElement()
	} else if g1.IsUnitElement() {
		gr.Set(g2)
	} else if g2.IsUnitElement() {
		gr.Set(g1)
	} else {
		lambdax := big.NewInt(0)
		lambday := big.NewInt(0)
		dx := big.NewInt(0)
		retx := big.NewInt(0)
		rety := big.NewInt(0)
		lambda := big.NewInt(0)

		fp.FpAdd(g1.Y, g2.Y, gr.Ec.P, rety)
		if rety.Cmp(zero) == 0 {
			gr.SetUnitElement()
		} else if g1.X.Cmp(g2.X) == 0 {
			gr.Double(g1)
		} else {
			fp.FpSub(g2.Y, g1.Y, gr.Ec.P, lambday)
			fp.FpSub(g2.X, g1.X, gr.Ec.P, lambdax)
			fp.FpDiv(lambday, lambdax, gr.Ec.P, lambda)

			fp.FpMul(lambda, lambda, gr.Ec.P, retx)
			fp.FpSub(retx, g1.X, gr.Ec.P, retx)
			fp.FpSub(retx, g2.X, gr.Ec.P, retx)

			fp.FpSub(g1.X, retx, gr.Ec.P, dx)
			fp.FpMul(lambda, dx, gr.Ec.P, rety)
			fp.FpSub(rety, g1.Y, gr.Ec.P, rety)

			gr.X.Set(retx)
			gr.Y.Set(rety)
		}
	}
	return gr
}

func (gr *ECElement) Double(g1 *ECElement) *ECElement {
	if g1.Y.CmpAbs(gr.Ec.P) == 0 {
		gr.X.Set(gr.Ec.P)
		gr.Y.Set(gr.Ec.P)
	} else {
		C := big.NewInt(0)
		lambda := big.NewInt(3)
		retx := big.NewInt(0)
		rety := big.NewInt(0)

		fp.FpMul(lambda, g1.X, gr.Ec.P, lambda)
		fp.FpMul(lambda, g1.X, gr.Ec.P, lambda)
		fp.FpAdd(lambda, gr.Ec.A, gr.Ec.P, lambda)
		fp.FpDiv(lambda, g1.Y, gr.Ec.P, lambda)
		fp.FpDiv(lambda, two, gr.Ec.P, lambda)

		retx.Set(lambda)
		fp.FpMul(retx, lambda, gr.Ec.P, retx)
		fp.FpSub(retx, g1.X, gr.Ec.P, retx)
		fp.FpSub(retx, g1.X, gr.Ec.P, retx)

		fp.FpSub(g1.X, retx, gr.Ec.P, C)
		fp.FpMul(lambda, C, gr.Ec.P, rety)
		fp.FpSub(rety, g1.Y, gr.Ec.P, rety)

		gr.X.Set(retx)
		gr.Y.Set(rety)
		//fmt.Println("Double", gr.Ec.P, lambda, retx, rety)
	}
	return gr
}

func (gr *ECElement) ScalarMul(g1 *ECElement, m *big.Int) *ECElement {
	if g1.Y.CmpAbs(gr.Ec.P) == 0 {
		gr.X.Set(gr.Ec.P)
		gr.Y.Set(gr.Ec.P)
	} else {
		// Set unit elemet to return value
		gret := NewECElement(gr.Ec, gr.Ec.P, gr.Ec.P)
		gtmp := NewECElement(gr.Ec, g1.X, g1.Y)

		for i := 0; i < m.BitLen(); i++ {
			if m.Bit(i) > 0 {
				gret.Add(gret, gtmp)
			}
			gtmp.Double(gtmp)
		}
		gr.Set(gret)
	}
	return gr
}

func GenerateSecp256k1PublicKey(secret *big.Int) (*ECElement, error) {
	g := NewECElement(&Secp256k1, Secp256k1.P, Secp256k1.P)

	if secret.Cmp(zero) == 0 || secret.Cmp(Secp256k1.Order) >= 0 {
		err := errors.New("invalid secret key")
		return g, err
	}

	gsecp := NewECElement(&Secp256k1, Secp256k1.Gx, Secp256k1.Gy)
	g.ScalarMul(gsecp, secret)
	return g, nil
}

/*
func GetSecp256k1SchnorrSignature(secret *big.Int, hash [32]byte) ([64]byte, error) {
	if secret.Cmp(zero) == 0 || secret.Cmp(Secp256k1.Order) >= 0 {
		err := errors.New("invalid secret key")
		return [64]byte{}, err
	}
	d := big.NewInt(0).Set(secret)

	p, err := GenerateSecp256k1PublicKey(d)
	if err != nil {
		return [64]byte{}, err
	}

	if p.Y.Bit(0) != 0 {
		fp.FpSub(Secp256k1.Order, p, Secp256k1.Order, d)
	}

	max, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)
	a, errrand := rand.Int(rand.Reader, max)

}
*/
