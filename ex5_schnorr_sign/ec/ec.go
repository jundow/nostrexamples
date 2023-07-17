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

func GetTaggedHash(tag string, items ...[]byte) [32]byte {
	tagb := []byte(tag)
	tagh := sha256.Sum256(tagb)
	h := sha256.New()
	h.Write(tagh[:])
	h.Write(tagh[:])
	for i := 0; i < len(items); i++ {
		h.Write(items[i])
	}

	var ret [32]byte
	copy(ret[0:32], h.Sum(nil))
	return ret
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
	var pub_b [32]byte
	copy(pub_b[0:32], p.X.FillBytes(make([]byte, 32)))

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
	rand_h := GetTaggedHash("BIP0340/aux", a[0:32])

	for i := 0; i < 32; i++ {
		t[i] = a[i] ^ rand_h[i]
	}

	//To the random number to sign as;
	//rand = hash(bytes("BIP0430/nonce") || bytes("BIP0430/nonce") || t || Pubkey || message )
	rand := GetTaggedHash("BIP0430/nonce", t[0:32], pub_b[0:32], message)

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
	//Elliptic curve points at r = kd g in byte array
	var r_b [32]byte
	copy(r_b[0:32], r.X.FillBytes(make([]byte, 32)))

	e_h := GetTaggedHash("BIP0340/challenge", r_b[0:32], pub_b[0:32], message)

	e := big.NewInt(0).SetBytes(e_h[0:32])
	e.Mod(e, Secp256k1.Order)

	//let s_b be bytes( k+ed mod order )
	s := big.NewInt(0).Set(e)
	s.Mul(s, d)
	s.Add(s, k)
	s.Mod(s, Secp256k1.Order)

	var s_b [32]byte
	copy(s_b[0:32], s.FillBytes(make([]byte, 32)))

	//return r_b = bytes(r) and s_b
	var ret [64]byte
	copy(ret[0:64], append(r_b[0:32], s_b[0:32]...))

	return ret, nil
}

func VerifySecp256k1(public [32]byte, message []byte, sig [64]byte) (bool, error) {
	// pv.x = int(public)
	// py = sqrt(pv.x^3 + 7 mod p) ... Select even value of two possible solutions

	pv := NewECElement(&Secp256k1, Secp256k1.P, Secp256k1.P)

	pvx := big.NewInt(0).SetBytes(public[:])
	if pvx.Cmp(Secp256k1.P) >= 0 {
		return false, fmt.Errorf("Invalid Public Key")
	}

	pvy := big.NewInt(0)
	pvr := big.NewInt(0)
	fp.FpMul(pvx, pvx, Secp256k1.P, pvr)
	fp.FpMul(pvr, pvx, Secp256k1.P, pvr)
	fp.FpAdd(pvr, big.NewInt(7), Secp256k1.P, pvr)
	fp.FpSecp256kSqrt(pvr, Secp256k1.P, pvy)

	pvy2 := big.NewInt(0)
	fp.FpMul(pvy, pvy, Secp256k1.P, pvy2)

	if pvr.Cmp(pvy2) != 0 {
		return false, fmt.Errorf("Invalid Public Key")
	}

	pv.X.Set(pvx)
	pv.Y.Set(pvy)

	// Derive r from sig[0:32]
	r := big.NewInt(0).SetBytes(sig[0:32])
	if r.Cmp(Secp256k1.P) >= 0 {
		return false, fmt.Errorf("Invalid signature sig[0:32]")
	}

	s := big.NewInt(0).SetBytes(sig[32:64])
	if s.Cmp(Secp256k1.Order) >= 0 {
		return false, fmt.Errorf("Invalid signature sig[32:64]")
	}

	//e = int( Hash(sha256(tag) || sha256(tag) || bytes(r) || bytes(pv.X) ) ) mod n
	var pvx_b [32]byte
	copy(pvx_b[0:32], pv.X.FillBytes(make([]byte, 32)))

	e_hash := GetTaggedHash("BIP0340/challenge", sig[0:32], pvx_b[0:32], message)

	e := big.NewInt(0).SetBytes(e_hash)
	e.Mod(e, Secp256k1.Order)

	//rv = sG - ePv

	ptmp := NewECElement(&Secp256k1, Secp256k1.P, Secp256k1.P)
	ptmp.ScalarMul(pv, e)
	//Negate
	ptmp.Y.Sub(Secp256k1.P, ptmp.Y)
	ptmp.Y.Mod(ptmp.Y, Secp256k1.P)

	rv := NewECElement(&Secp256k1, Secp256k1.Gx, Secp256k1.Gy)
	rv.ScalarMul(rv, s)
	rv.Add(rv, ptmp)

	//Check the signature is valid and retun the result

	if rv.X.Cmp(Secp256k1.P) == 0 {
		return false, fmt.Errorf("Invalid signature: rv is at infinity")
	}

	if rv.Y.Bit(0) != 0 {
		return false, fmt.Errorf("Invalid signature: rv is not even")
	}

	if rv.X.Cmp(r) != 0 {
		return false, fmt.Errorf("Invalid signature: rv.X does not match to r")
	}

	return true, nil

}
