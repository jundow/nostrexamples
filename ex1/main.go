package main

import (
	"encoding/json"
	"fmt"

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

		//fmt.Println(string(rmsg))

		if msgtyp := (dat.([]any))[0]; msgtyp == "EOSE" {
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

func main() {
	authors := "[\"" +
		"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX" + //Hex pubkey comes here
		"\"]"
	urlws := "wss://nos.lol/"
	urlhttp := "https://nos.lol/"

	msg := "[\"REQ\", \"1\", {\"kinds\":[1], \"authors\":" + authors + ",\"limit\":100}]"
	var v []any

	ws, wserr := websocket.Dial(urlws, "", urlhttp)
	if wserr != nil {
		fmt.Println(wserr)
		return
	}

	//fmt.Println(msg)

	Send(ws, msg)
	rerr := Recv(ws, &v)
	defer ws.Close()

	if rerr != nil {
		fmt.Println("Read Error.")
		fmt.Println(rerr)
		return
	}

	for _, item := range v {
		fmt.Println(item)
	}
}
