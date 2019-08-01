package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/ext/bridgemon"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

const audiosocketEndpoint = "AudioSocket/%s:8080/%s"

// LocalChannelAnswerTimeout is the maximum time to wait for a local channel to be answered
var LocalChannelAnswerTimeout = time.Second

// State is the structure for storing application execution data
type State struct {
	h *ari.ChannelHandle

	tries int
}

type stateFn func(context.Context) (stateFn, error)

func init() {
	if os.Getenv("AUDIOSOCKET_SERVICE_HOST") == "" {
		os.Setenv("AUDIOSOCKET_SERVICE_HOST", "localhost")
	}
}

func app(ctx context.Context, ac ari.Client, h *ari.ChannelHandle) error {
	log.Println("running voice app")

	// Always quit on hangup
	go func() {
		sub := h.Subscribe(ari.Events.ChannelDestroyed)
		defer sub.Cancel()
		select {
		case <-sub.Events():
		case <-ctx.Done():
		}
	}()

	h.Answer()
	time.Sleep(time.Second)

	id := uuid.Must(uuid.NewV1())

	// Bridge to voice app (via AudioSocket)
	br, err := ac.Bridge().Create(h.Key(), "mixing", "bridge-"+id.String())
	if err != nil {
		return errors.Wrap(err, "failed to stage bridge creation")
	}
	m := bridgemon.New(br)
	defer m.Close()

	brEvents := m.Watch()

	if err := br.AddChannel(h.ID()); err != nil {
		return errors.Wrap(err, "failed to add original channel to bridge")
	}
	defer br.RemoveChannel(h.ID()) // nolint

	// Create an AudioSocket channel
	as, err := ac.Channel().StageOriginate(h.Key(), ari.OriginateRequest{
		Endpoint:   fmt.Sprintf(audiosocketEndpoint, os.Getenv("AUDIOSOCKET_SERVICE_HOST"), id.String()),
		ChannelID:  id.String(),
		App:        ac.ApplicationName(),
		AppArgs:    "noop",
		Originator: h.ID(),
		Variables: map[string]string{
			"AUDIOSOCKET_ID": id.String(),
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to stage AudioSocket channel creation")
	}

	asStart := as.Subscribe(ari.Events.StasisStart)
	go func() {
		defer asStart.Cancel()

		select {
		case <-ctx.Done():
			return
		case <-time.After(LocalChannelAnswerTimeout):
			return
		case <-asStart.Events():
			break
		}

		// Send AudioSocket channel to the bridge
		if err := br.AddChannel(as.ID()); err != nil {
			log.Println("failed to send AudioSocket channel to bridge:", err)
			return
		}
	}()

	if err := as.Exec(); err != nil {
		return errors.Wrap(err, "failed to create AudioSocket channel")
	}
	defer as.Hangup() // nolint: errcheck

	// Wait for the bridge to be left
	log.Println("waiting for bridge quorum")
	var hadQuorum bool
	for {
		select {
		case <-ctx.Done():
			log.Println("context terminated")
			return nil
		case data := <-brEvents:
			if len(data.ChannelIDs) > 1 {
				log.Println("bridge quorum achieved")
				hadQuorum = true
			}
			if len(data.ChannelIDs) < 2 && hadQuorum {
				log.Println("channel left bridge; exiting")
				return nil
			}
			log.Printf("odd bridge state with %d members", len(data.ChannelIDs))
		}
	}
}
