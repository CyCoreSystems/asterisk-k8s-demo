package main

import (
	"log"
	"time"

	"golang.org/x/net/context"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/ext/audio"
	"github.com/CyCoreSystems/ari/ext/prompt"
)

// State is the structure for storing application execution data
type State struct {
	ctx context.Context

	h *ari.ChannelHandle

	digit string
}

type stateFn func(*State) stateFn

func app(cl *ari.Client, h *ari.ChannelHandle) {
	log.Println("Running channel app")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Minute))

	// Always quit on hangup
	go func() {
		sub := h.Subscribe(ari.Events.ChannelDestroyed)
		select {
		case <-sub.Events():
			cancel()
		case <-ctx.Done():
		}
		sub.Cancel()
	}()

	h.Answer()

	// Create the state struct
	s := &State{
		ctx: ctx,
		h:   h,
	}

	// Run state machine
	for next := menu; next != nil; {
		next = next(s)
	}

	h.Hangup()

}

func menu(s *State) stateFn {
	ret, err := prompt.Prompt(s.ctx, s.h, nil, "sound:hello", "sound:please-enter-your", "sound:number")
	if err != nil {
		return nil
	}

	if ret.Status != prompt.Complete {
		return invalid
	}

	s.digit = ret.Data
	return reply
}

func reply(s *State) stateFn {
	q := audio.NewQueue()
	q.Add("sound:you-entered")
	q.Add("digits:" + s.digit)
	_, _ = q.Play(s.ctx, s.h, nil)
	return nil
}

func invalid(s *State) stateFn {
	audio.Play(s.ctx, s.h, "an-error-has-occurred")
	return nil
}
