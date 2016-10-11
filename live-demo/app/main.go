package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/client"
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

	cl, err := connect(ctx)
	if err != nil {
		log.Println("Failed to build nc ARI client", "error", err)
		return -1
	}

	// setup app

	log.Println("Starting listener app")

	listen(cl, app)

	return 0
}

func listen(cl *ari.Client, handler func(cl *ari.Client, h *ari.ChannelHandle)) {
	sub := cl.Bus.Subscribe("StasisStart")

	for e := range sub.Events() {
		log.Println("Got stasis start")
		stasisStartEvent := e.(*ari.StasisStart)
		go handler(cl, cl.Channel.Get(stasisStartEvent.Channel.ID))
	}
}

func connect(ctx context.Context) (cl *ari.Client, err error) {

	opts := client.Options{
		URL: "nats://nats:4222",
	}

	log.Println("Connecting")

	cl, err = client.New(opts)
	return
}
