package main

import (
	"context"
	"log"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/client"
)

const ariApp = "demo"

var baseClient *client.Client

func main() {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect
	log.Println("connecting to ARI")
	baseClient, err = client.New(ctx, client.WithApplication(ariApp))
	if err != nil {
		log.Println("failed to build ARI client", "error", err)
		return
	}

	log.Println("starting listener")
	err = client.Listen(ctx, baseClient, appStart)
	if err != nil {
		log.Println("failed to listen for new calls")
	}
	<-ctx.Done()

	return
}

func appStart(h *ari.ChannelHandle, startEvent *ari.StasisStart) {
	log.Println("running app:", "channel", h.Key().ID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Minute))
	defer cancel()

	if err := app(ctx, baseClient.New(ctx), h); err != nil {
		log.Println("app execution failed:", err.Error())
	}

	h.Hangup()
	log.Println("channel hung up")
	return
}
