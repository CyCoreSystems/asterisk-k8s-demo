package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/client"
	"github.com/CyCoreSystems/ari-proxy/session"
	"github.com/nats-io/nats"
)

func main() {

	<-time.After(1 * time.Second)

	if i := run(); i != 0 {
		os.Exit(i)
	}
}

func run() int {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect

	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Println("Failed to connect to NATS", "error", err)
		return -1
	}

	// setup app

	log.Println("Starting listener app")

	err = client.Listen(ctx, nc, "demo", handler)
	if err != nil {
		log.Println("Unable to start ari-proxy client listener:", err)
		return -2
	}

	return 0
}

func handler(cl *ari.Client, d *session.Dialog) {
	log.Println("Starting dialog handler", d.ID)
	h := cl.Channel.Get(d.ChannelID)
	app(cl, h)
}
