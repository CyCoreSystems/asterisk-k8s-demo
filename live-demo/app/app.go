package main

import (
	"log"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/ext/audio"
	"github.com/CyCoreSystems/ari/ext/prompt"
	"github.com/CyCoreSystems/dispatchers/deployment"
)

// State is the structure for storing application execution data
type State struct {
	ctx context.Context

	h *ari.ChannelHandle

	digit string

	tries int
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
	time.Sleep(time.Second)

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
	if s.tries > 3 {
		return invalid
	}
	s.tries++

	ret, err := prompt.Prompt(s.ctx, s.h, nil, "sound:hello", "sound:please-enter-your", "sound:number")
	if err != nil {
		log.Println("Failed to play prompt", err)
		return invalid
	}

	switch ret.Status {
	case prompt.Complete:
		s.digit = ret.Data
		return reply
	case prompt.Timeout:
		return menu
	default:
		log.Println("Prompt status not complete", ret.Status)
		return invalid
	}
}

func reply(s *State) stateFn {
	log.Println("Sending reply", s.digit)
	q := audio.NewQueue()
	q.Add("sound:you-entered")
	q.Add("digits:" + s.digit)
	_, _ = q.Play(s.ctx, s.h, nil)
	return scale
}

func scale(s *State) stateFn {
	ni, err := strconv.Atoi(s.digit)
	if err != nil {
		log.Println("Failed to parse digit as number", err)
		return nil
	}
	n := int32(ni)

	err = deployment.Scale(&n)
	if err != nil {
		log.Println("Failed to scale asterisk", err)
		return nil
	}
	log.Println("Scaled asterisk to ", n)
	return nil
}

func invalid(s *State) stateFn {
	log.Println("Playing error message")
	audio.Play(s.ctx, s.h, "sound:an-error-has-occurred")
	return nil
}
