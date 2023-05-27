package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
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

func Sign(seck string, hash [32]byte) ([64]byte, error) {

	sekbyte, sekerr := hex.DecodeString(seck)
	if sekerr != nil {
		return [64]byte{}, sekerr
	}

	sek, _ := btcec.PrivKeyFromBytes(sekbyte)
	//pubkbyte := pubk.SerializeCompressed()
	//pubkhex := hex.EncodeToString(pubkbyte[1:])
	// I am not sure what is happening but this will create the same key as pubkey
	evSig, sigerr := schnorr.Sign(sek, hash[:])

	if sigerr != nil {
		return [64]byte{}, sigerr
	}
	return [64]byte(evSig.Serialize()), nil
}

func main() {
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
