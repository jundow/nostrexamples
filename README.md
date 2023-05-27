## Examples of minimally coded Nostr program

These examples aim to understand how Nostr relay works with NIPs.

### ex1 (Golang)
Reads single user's notes from single relay.

Hardcode a hex public key and the URL of a relay both in "ws" and "http" endpoint.
The corresponding code is at the top of "main" function.

``` golang
authors := "[\"" +
	"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX" + //Hex pubkey comes here
	"\"]"
urlws := "wss://nos.lol/"
urlhttp := "https://nos.lol/"
```

and

``` bash
$ go run main.go
```

to get the recent 100 notes of the corresponding user. 

### ex2 (Golang)
Reads single user's follows and obtain notes from them from a single relay.

You only need to hardcode a public key and a relay's URLs into the source code,

``` golang
urlws := "wss://nos.lol/"
urlhttp := "https://nos.lol/"
pub := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
```

and

``` bash
$ go run main.go
```

to get the recent 1000 notes from the time line.

## ex3 (Golang)

Probably the most minimal program that post a note to single relay.

However, I could not create an elliptic curve cryptograph program from the scrathc.
You need to go get elliptic cure cryptograph library as follows.

``` bash
$ go get github.com/btcsuite/btcd/btcec/v2/schnorr
```

**CAUTION** 
You need to hardcode urls of relay, public key and **secret key.**

It is strongly recommended to prepare a key-pair other than the keys associated with your Nostr identity to avoid an accidenttal key leakage.

``` golang
	urlws := "wss://nos.lol/"
	urlhttp := "https://nos.lol/"
	pubkhex := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	sekhex := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	msg := "Your message comes here."
```
Put your message to be issued to the variable "msg."
The event will immediately be issued once the program ran.

