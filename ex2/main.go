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

func GetFollows(pub string, relay string, relayhttp string) ([]string, error) {
	ws, wserr := websocket.Dial(relay, "", relayhttp)
	if wserr != nil {
		return nil, wserr
	}

	var v []any
	msg := "[\"REQ\"," +
		"\"0\"," +
		"{\"kinds\": [3]," +
		"\"authors\": [\"" + pub + "\"]}]"
	Send(ws, msg)
	Recv(ws, &v)
	defer ws.Close()

	var follows []string

	for _, item := range v[0].([]any)[2].(map[string]any)["tags"].([]any) {
		if item.([]any)[0].(string) == "p" {
			follows = append(follows, item.([]any)[1].(string))
			//fmt.Println(item.([]any)[1].(string))
		}
	}
	return follows, nil
}

func main() {
	urlws := "wss://nos.lol/"
	urlhttp := "https://nos.lol/"
	pub := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	follows, err := GetFollows(pub, urlws, urlhttp)

	if err != nil {
		fmt.Println(err)
		return
	}

	var authors string
	authors = "["
	for _, follow := range follows {
		authors += "\"" + follow + "\","
	}
	authors = authors[:len(authors)-1] + "]"

	msg := "[\"REQ\", \"1\", {\"kinds\": [1], \"authors\": " + authors + ",\"limit\": 1000}]"

	var v []any
	ws, wserr := websocket.Dial(urlws, "", urlhttp)
	if wserr != nil {
		fmt.Println(wserr)
		return
	}
	Send(ws, msg)
	Recv(ws, &v)
	defer ws.Close()

	fmt.Println(msg)
	for _, item := range v {
		fmt.Println(item)
	}
}
