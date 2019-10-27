package main

import (
	"fmt"
	"log"

	"github.com/CyCoreSystems/agi"
	"github.com/gofrs/uuid"
)

const listenAddr = ":8080"
const audiosocketAddr = "audiosocket:8080"

func main() {
	log.Fatalln(agi.Listen(listenAddr, callHandler))
}

func callHandler(a *agi.AGI) {
	log.Println("new call from", a.Variables["agi_callerid"])

	id := uuid.Must(uuid.NewV1())

	defer func() {
		a.Hangup() // nolint
		a.Close()  // nolint
	}()

	if err := a.Answer(); err != nil {
		log.Println("failed to answer call:", err)
		return
	}

	if _, err := a.Exec("AudioSocket", fmt.Sprintf("%s,%s", id.String(), audiosocketAddr)); err != nil {
		log.Printf("failed to execute AudioSocket to %s: %v", audiosocketAddr, err)
	}
}
