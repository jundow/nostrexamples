package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"golang.org/x/net/websocket"

	"eccalc/ec"
	"eccalc/fp"
)

func Recv(ws *websocket.Conn, v *[]any) error {
	var rmsg []byte
	for {
		var dat any
		werr := websocket.Message.Receive(ws, &rmsg)
		if werr != nil {
			return werr
		}

		jerr := json.Unmarshal(rmsg, &dat)
		if jerr != nil {
			return jerr
		}

		fmt.Println(string(rmsg))

		if msgtyp := (dat.([]any))[0]; msgtyp == "EOSE" {
			return nil
		} else if msgtyp == "OK" {
			return nil
		} else if msgtyp == "NOTICE" {
			return fmt.Errorf("notice: %s", dat)
		} else {
			*v = append(*v, dat)
		}
	}
}

func Send(ws *websocket.Conn, msg string) {
	websocket.Message.Send(ws, msg)
}

func Serialize(pub string, created_at int64, kind int, tags [][]string, content string) [32]byte {
	str := "[0," +
		"\"" + pub + "\"," +
		fmt.Sprint(created_at) + "," +
		fmt.Sprint(kind) + "," +
		fmt.Sprint(tags) + "," +
		"\"" + content + "\"]"

	hash := sha256.Sum256([]byte(str))

	//fmt.Println(str)
	//fmt.Println([]byte(str))
	//fmt.Println(hash)

	return hash
}

func main() {
	/*
		z := big.NewInt(0)
		ret := big.NewInt(0)

		p64 := int64(31)
		p := big.NewInt(p64)

		for i := int64(0); i < p64; i++ {
			x := big.NewInt(i)
			for j := int64(0); j < p64; j++ {
				y := big.NewInt(j)
				fp.FpMul(x, y, p, z)
				fmt.Print(z)
				fmt.Print("\t")
			}
			fmt.Println()
		}

		for i := int64(1); i < p64; i++ {
			x := big.NewInt(i)
			fp.FpInv(x, p, z)
			fmt.Println(x, z, fp.FpMul(x, z, p, ret))
		}

		fmt.Println("######")

		pow := big.NewInt(0)
		m := big.NewInt(3)
		for i := int64(0); i < p64; i++ {
			pow.SetInt64(i)
			fp.FpPow(m, pow, p, ret)
			fmt.Println("ANSWER", m, pow, p, ret)
		}

		pSecp256k, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)

		for i := int64(2); i < p64; i++ {
			x := big.NewInt(i)
			_, err := fp.FpSecp256kSqrt(x, pSecp256k, z)
			fmt.Println(err, x, z)
		}

		gs := ec.NewECElement(&ec.Test, ec.Test.Gx, ec.Test.Gy)
		g := ec.NewECElement(&ec.Test, ec.Test.Gx, ec.Test.Gy)

		var i int64
		i = 1

		for i = 1; i < 100; i++ {
			g.ScalarMul(gs, m.SetInt64(i))
			fmt.Println(i, g.X, g.Y)
		}
	*/

	/*

	filep, err := os.Open("../../testkeys")
	if err != nil {
		return
	}
	defer filep.Close()

	skeys := []string{}
	scanner := bufio.NewScanner(filep)
	for scanner.Scan() {
		line := scanner.Text()
		skeys = append(skeys, line)
	}

	s := big.NewInt(0)
	gsecp := ec.NewECElement(&ec.Secp256k1, ec.Secp256k1.Gx, ec.Secp256k1.Gy)

	var g *ec.ECElement

	for _, skey := range skeys {
		seckey, _ := big.NewInt(0).SetString(skey, 16)
		g = ec.NewECElement(&ec.Secp256k1, ec.Secp256k1.Gx, ec.Secp256k1.Gy)
		g.ScalarMul(gsecp, seckey)
		s.Set(g.X)
		fp.FpPow(s, big.NewInt(3), g.Ec.P, s)
		fp.FpAdd(s, big.NewInt(7), g.Ec.P, s)
		fp.FpSecp256kSqrt(s, g.Ec.P, s)

		fmt.Println(fmt.Sprintf("skey %x", seckey))
		fmt.Println(fmt.Sprintf("X    %x", g.X))
		fmt.Println(fmt.Sprintf("Y    %x", g.Y))
		fmt.Println(fmt.Sprintf("S    %x", s))

		if s.Bit(0) != 0 {
			fmt.Println("odd")
			fp.FpSub(g.Ec.P, s, g.Ec.P, s)
			fmt.Println(fmt.Sprintf("S    %x", s))
		}
		fmt.Println()
	}
	*/

	filep, err := os.Open("../../testkeys")
	if err != nil {
		return
	}
	defer filep.Close()

	skeys := []([32]byte)

	skeystr := []string{}
	scanner := bufio.NewScanner(filep)
	for scanner.Scan() {
		line := scanner.Text()
		skeystr = append(skeystr, line)


	}

	urlws := "wss://nos.lol/"
	urlhttp := "https://nos.lol/"
	pubkhex := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	sekhex := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	msg := "Your message comes here."

	tags := [][]string{}
	kind := 1
	created_at := time.Now().Unix()

	hash := Serialize(pubkhex, created_at, kind, tags, msg)
	evIDhex := hex.EncodeToString(hash[:])
	evSig, sigerr := Sign(sekhex, hash)
	if sigerr != nil {
		fmt.Println("Signature Error")
		fmt.Println(sigerr)
	}
	evSighex := hex.EncodeToString(evSig[:])

	eventstr := "[\"EVENT\",{" +
		"\"id\":" + "\"" + evIDhex + "\"," +
		"\"pubkey\":" + "\"" + pubkhex + "\"," +
		"\"created_at\":" + fmt.Sprint(created_at) + "," +
		"\"kind\":" + fmt.Sprint(kind) + "," +
		"\"tags\":" + fmt.Sprint(tags) + "," +
		"\"content\":" + "\"" + msg + "\"," +
		"\"sig\":" + "\"" + evSighex + "\"" +
		"}]"

	fmt.Println(eventstr)

	var v []any
	ws, wserr := websocket.Dial(urlws, "", urlhttp)
	if wserr != nil {
		fmt.Println(wserr)
		return
	}
	Send(ws, eventstr)
	Recv(ws, &v)
	defer ws.Close()

	fmt.Println(msg)
	for _, item := range v {
		fmt.Println(item)
	}
}
