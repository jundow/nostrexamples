package main

import (
	"bufio"
	"crypto/sha256"
	"eccalc/ec"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/net/websocket"
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

func Serialize(pubkey [32]byte, created_at int64, kind int, tags [][]string, content string) [32]byte {

	pubstr := hex.EncodeToString(pubkey[0:32])

	str := "[0," +
		"\"" + pubstr + "\"," +
		fmt.Sprint(created_at) + "," +
		fmt.Sprint(kind) + "," +
		fmt.Sprint(tags) + "," +
		"\"" + content + "\"]"

	hash := sha256.Sum256([]byte(str))

	return hash
}

func main() {

	filep, err := os.Open("../../mkey")
	if err != nil {
		return
	}
	defer filep.Close()

	var skeys []([32]byte)
	var pkeys []([32]byte)

	scanner := bufio.NewScanner(filep)
	for scanner.Scan() {
		line := scanner.Text()
		tmpkey, skeyerr := hex.DecodeString(line)
		if skeyerr != nil {
			fmt.Println("Error hex.DecoteString ", line)
			return
		}
		var skey [32]byte
		copy(skey[0:32], tmpkey)
		skeys = append(skeys, skey)
		pkey, pkerr := ec.GenerateSecp256k1PublicKey(skey)
		if pkerr != nil {
			fmt.Println("Error Public Key Generation ")
			return
		}
		//fmt.Println(line)
		//fmt.Println(pkey)
		pkeys = append(pkeys, pkey)
	}

	msg := "テスト01 2023-07-22"

	tags := [][]string{}
	kind := 1
	created_at := time.Now().Unix()

	for i := 0; i < len(skeys); i++ {
		hash := Serialize(pkeys[i], created_at, kind, tags, msg)
		evIDhex := hex.EncodeToString(hash[:])
		//evSig, sigerr := Sign(sekhex, hash)
		evSig, sigerr := ec.SingSecp256k1(skeys[i], hash[:])
		if sigerr != nil {
			fmt.Println("Signature error", sigerr)
			return
		}
		evSighex := hex.EncodeToString(evSig[:])

		vres, verr := ec.VerifySecp256k1(pkeys[i], hash[:], evSig)

		if !vres {
			fmt.Println("Signature Verification error", verr)
			return
		}

		eventstr := "[\"EVENT\",{" +
			"\"id\":" + "\"" + evIDhex + "\"," +
			"\"pubkey\":" + "\"" + hex.EncodeToString(pkeys[i][0:32]) + "\"," +
			"\"created_at\":" + fmt.Sprint(created_at) + "," +
			"\"kind\":" + fmt.Sprint(kind) + "," +
			"\"tags\":" + fmt.Sprint(tags) + "," +
			"\"content\":" + "\"" + msg + "\"," +
			"\"sig\":" + "\"" + evSighex + "\"" +
			"}]"

		fmt.Println(eventstr)

		relays := []string{
			"nos.lol/",
			"relay.nostr.wirednet.jp",
			//"nostr.h3z.jp",
			//"nostr-world.h3z.jp",
			"nostr-relay.nokotaro.com",
		}

		for j := 0; j < len(relays); j++ {

			var v []any
			ws, wserr := websocket.Dial("wss://"+relays[j], "", "https://"+relays[j])
			if wserr != nil {
				fmt.Println(wserr)
				return
			}
			Send(ws, eventstr)
			Recv(ws, &v)
			defer ws.Close()

			fmt.Println(relays[j])
			fmt.Println(msg)
			for _, item := range v {
				fmt.Println(item)
			}
		}
	}
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
}
