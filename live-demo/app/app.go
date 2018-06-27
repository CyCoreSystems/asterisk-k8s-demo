package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/ext/play"
	"github.com/ericchiang/k8s"
	"github.com/ericchiang/k8s/apis/apps/v1"
	"github.com/pkg/errors"
)

// State is the structure for storing application execution data
type State struct {
	h *ari.ChannelHandle

	tries int
}

type stateFn func(context.Context) (stateFn, error)

func app(ctx context.Context, h *ari.ChannelHandle) error {
	log.Println("running channel app")

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

	// Create the state struct
	s := &State{
		h: h,
	}

	// Run state machine
	var err error
	for next := s.menu; next != nil; {
		if next, err = next(ctx); err != nil {
			if iErr := invalid(ctx, h); iErr != nil {
				log.Println("failed to play invalid message:", iErr)
			}

			return err
		}
	}
	return nil
}

func (s *State) menu(ctx context.Context) (stateFn, error) {
	if s.tries > 3 {
		return nil, errors.New("too many retries")
	}
	s.tries++

	ret, err := play.Prompt(ctx, s.h, play.URI("sound:hello", "sound:please-enter-your", "sound:number")).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to play prompt")
	}

	switch ret.MatchResult {
	case play.Complete:
	case play.Incomplete:
		return s.menu, nil
	default:
		return nil, errors.New("invalid DTMF entry")
	}

	size, err := strconv.Atoi(ret.DTMF)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert DTMF entry to an integer")
	}
	return s.reply(size), nil
}

func (s *State) reply(size int) func(context.Context) (stateFn, error) {
	return func(ctx context.Context) (stateFn, error) {
		log.Println("announcing new size:", size)
		err := play.Play(ctx, s.h, play.URI("sound:you-entered", fmt.Sprintf("digits:%d", size))).Err()
		return s.scale(size), err
	}
}

func (s *State) scale(size int) func(context.Context) (stateFn, error) {
	return func(ctx context.Context) (stateFn, error) {

		if err := scaleAsterisk(ctx, int32(size)); err != nil {
			return nil, errors.Wrap(err, "failed to scale asterisk deployment")
		}
		log.Println("scaled asterisk to ", size, "replicas")
		return nil, nil
	}
}

func invalid(ctx context.Context, h *ari.ChannelHandle) error {
	return play.Play(ctx, h, play.URI("sound:an-error-has-occurred")).Err()
}

func scaleAsterisk(ctx context.Context, size int32) error {
	k, err := k8s.NewInClusterClient()
	if err != nil {
		return errors.Wrap(err, "failed to get kubernetes client")
	}

	d := new(v1.Deployment)
	if err = k.Get(ctx, "voip", "asterisk", d); err != nil {
		return errors.Wrap(err, "failed to retrieve current asterisk deployment")
	}

	current := d.GetSpec().GetReplicas()
	if current == size {
		return nil
	}

	d.GetSpec().Replicas = &size
	return k.Update(ctx, d)
}
