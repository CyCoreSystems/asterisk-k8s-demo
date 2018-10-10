package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/ext/bridgemon"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

const audiosocketEndpoint = "Local/audiosocket@default/n"

// LocalChannelAnswerTimeout is the maximum time to wait for a local channel to be answered
var LocalChannelAnswerTimeout = time.Second

// State is the structure for storing application execution data
type State struct {
	h *ari.ChannelHandle

	tries int
}

type stateFn func(context.Context) (stateFn, error)

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
	br, err := ac.Bridge().StageCreate(h.Key(), "mixing", id.String())
	if err != nil {
		return errors.Wrap(err, "failed to stage bridge creation")
	}
	m := bridgemon.New(br)
	defer m.Close()

	brEvents := m.Watch()

	if err := br.Exec(); err != nil {
		return errors.Wrap(err, "failed to create bridge")
	}
	legAID := "a-" + id.String()
	legBID := "b-" + id.String()

	// Create local channel pair for AudioSocket connection
	legA, err := ac.Channel().StageOriginate(h.Key(), ari.OriginateRequest{
		Endpoint:       audiosocketEndpoint,
		ChannelID:      legAID,
		OtherChannelID: legBID,
		App:            ac.ApplicationName(),
		AppArgs:        "noop",
		Originator:     h.ID(),
		Variables: map[string]string{
			"AUDIOSOCKET_ID": id.String(),
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to stage local channel creation")
	}
	legB := ari.NewChannelHandle(legA.Key().New(ari.ChannelKey, legBID), ac.Channel(), nil)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		legAStart := legA.Subscribe(ari.Events.StasisStart)
		defer legAStart.Cancel()
		wg.Done()

		select {
		case <-ctx.Done():
			return
		case <-time.After(LocalChannelAnswerTimeout):
			return
		case <-legAStart.Events():
			break
		}

		// Send LegA to the bridge
		if err := br.AddChannel(legAID); err != nil {
			log.Println("failed to send legA to brige:", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		legBStart := legB.Subscribe(ari.Events.StasisStart)
		defer legBStart.Cancel()

		wg.Done()

		select {
		case <-ctx.Done():
			return
		case <-time.After(LocalChannelAnswerTimeout):
			return
		case <-legBStart.Events():
		}

		// Send LegB to the AudioSocket
		if err := legB.Continue("inbound", "audiosocket", 1); err != nil {
			log.Println("failed to send legB to AudioSocket:", err)
			return
		}
	}()

	if err := legA.Exec(); err != nil {
		return errors.Wrap(err, "failed to create local channel pair")
	}

	// Add our original channel to the bridge
	if err := br.AddChannel(h.ID()); err != nil {
		return errors.Wrap(err, "failed to add original channel to bridge")
	}
	defer br.RemoveChannel(h.ID()) // nolint

	// Wait for the bridge to be left
	var hadQuorum bool
	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-brEvents:
			if len(data.ChannelIDs) > 1 {
				hadQuorum = true
			}
			if len(data.ChannelIDs) < 2 && hadQuorum {
				return nil
			}
		}
	}
}
