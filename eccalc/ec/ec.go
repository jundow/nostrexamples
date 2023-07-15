package ec

import (
	"crypto/rand"
	"crypto/sha256"
	"eccalc/fp"
	"errors"
	"fmt"
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
var rand_max *big.Int

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
	rand_max, _ = big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)
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

func GenerateSecp256k1PublicKey(sec [32]byte) ([32]byte, error) {

	secret := big.NewInt(0).SetBytes(sec[0:32])

	g := NewECElement(&Secp256k1, Secp256k1.P, Secp256k1.P)

	if secret.Cmp(zero) == 0 || secret.Cmp(Secp256k1.Order) >= 0 {
		err := errors.New("invalid secret key")
		return [32]byte{}, err
	}

	gsecp := NewECElement(&Secp256k1, Secp256k1.Gx, Secp256k1.Gy)
	g.ScalarMul(gsecp, secret)

	var ret [32]byte
	tmp := make([]byte, 32)

	g.X.FillBytes(tmp)
	copy(ret[0:32], tmp[0:32])

	return ret, nil
}

func (gr *ECElement) GetBytes() ([32]byte, [32]byte) {
	var bufx, bufy [32]byte
	gr.X.FillBytes(bufx[0:32])
	gr.Y.FillBytes(bufy[0:32])
	return bufx, bufy
}

func SingSecp256k1(secret [32]byte, message []byte) ([64]byte, error) {

	//Generate a random byte array as rand
	raux, raux_err := rand.Int(rand.Reader, rand_max)
	if raux_err != nil {
		return [64]byte{}, raux_err
	}
	var a [32]byte
	copy(a[0:32], raux.FillBytes(make([]byte, 32)))

	////////////////////////
	//Public Key Generation

	//sec := int(secret)
	//Convert the byte array secret key to big.int

	dd := big.NewInt(0).SetBytes(secret[0:32])

	if dd.Cmp(Secp256k1.Order) >= 0 {
		return [64]byte{}, fmt.Errorf("Invalid Secret Key")
	}

	//p = d'G, G as the generator of secp256k1 elliptic curve
	//Calculate the public key as the point of secp256k1 elliptic curve
	gsecp := NewECElement(&Secp256k1, Secp256k1.Gx, Secp256k1.Gy)
	p := NewECElement(&Secp256k1, Secp256k1.Gx, Secp256k1.Gy)
	p.ScalarMul(gsecp, dd)

	//public key in byte array
	var bytes_Pub [32]byte
	copy(bytes_Pub[0:32], p.X.FillBytes(make([]byte, 32)))

	//If p.Y is even let d = dd, otherwise d = order -d (mod Order of secp256k1)
	var d *big.Int
	if p.Y.Bit(0) == 0 {
		d = big.NewInt(0).Set(dd)
	} else {
		d = big.NewInt(0)
		d.Sub(Secp256k1.Order, dd)
	}

	/////////////////////////////////
	//Nonce(random number) generation

	//Let t be the xor of d and hash ( bytes("BIP0340/aux") || bytes("BIP0340/aux") || rand )
	// xor will be calculated each byte,a is the byte arrar of rand(big.int)
	var t [32]byte
	tagBIP0340aux := []byte("BIP0340/aux")
	rand_h := sha256.Sum256(append(append(tagBIP0340aux, tagBIP0340aux...), a[0:32]...))
	for i := 0; i < 32; i++ {
		t[i] = a[i] ^ rand_h[i]
	}

	//To the random number to sign as;
	//rand = hash(bytes("BIP0430/nonce") || bytes("BIP0430/nonce") || t || Pubkey || message )
	tagBIP0340nonce := []byte("BIP0340/nonce")
	rand_b := append(tagBIP0340nonce, tagBIP0340nonce...)
	rand_b = append(rand_b, t[0:32]...)
	rand_b = append(rand_b, bytes_Pub[0:32]...)
	rand_b = append(rand_b, message...)

	rand := sha256.Sum256(rand_b[:])

	//k' = int(rand) mod n
	kd := big.NewInt(0).SetBytes(rand[0:32])
	kd.Mod(kd, Secp256k1.Order)

	//Fail kd == 0
	if kd.Cmp(zero) == 0 {
		return [64]byte{}, fmt.Errorf("Invalid K-dash value")
	}

	//r = kd g
	gsecp = NewECElement(&Secp256k1, Secp256k1.Gx, Secp256k1.Gy)
	r := NewECElement(&Secp256k1, Secp256k1.Gx, Secp256k1.Gy)
	r.ScalarMul(gsecp, kd)

	//let k = kd if r.Y is even, otherwise k = order - kd
	var k *big.Int
	if r.Y.Bit(0) == 0 {
		k = big.NewInt(0).Set(kd)
	} else {
		k = big.NewInt(0)
		k.Sub(Secp256k1.Order, kd)
	}

	//////////////////////
	//Generate a signature

	// let e the integer of hash( bytes("bytes/challenge") || bytes("bytes/challenge") || bytes(R) || bytes(P) || m )
	tagBIP0340challenge := []byte("BIP0340/challenge")
	tbc := sha256.Sum256(tagBIP0340challenge)

	//Elliptic curve points at r = kd g in byte array
	var r_b [32]byte
	copy(r_b[0:32], r.X.FillBytes(make([]byte, 32)))

	e_htmp := sha256.New()
	//e_htmp.Write(tagBIP0340challenge)
	//e_htmp.Write(tagBIP0340challenge)
	e_htmp.Write(tbc[0:32])
	e_htmp.Write(tbc[0:32])
	e_htmp.Write(r_b[0:32])
	e_htmp.Write(bytes_Pub[0:32])
	e_htmp.Write(message)
	e_h := e_htmp.Sum(nil)

	e := big.NewInt(0).SetBytes(e_h[0:32])
	e.Mod(e, Secp256k1.Order)

	//let s_b be bytes( k+ed mod order )
	s := big.NewInt(0).Set(e)
	s.Mul(s, d)
	s.Add(s, k)
	s.Mod(s, Secp256k1.Order)

	var s_b [32]byte
	sbtmp := make([]byte, 32)
	s.FillBytes(sbtmp)
	copy(s_b[0:32], sbtmp)

	//return r_b = bytes(r) and s_b
	var ret [64]byte
	copy(ret[0:64], append(r_b[0:32], s_b[0:32]...))

	return ret, nil
}
