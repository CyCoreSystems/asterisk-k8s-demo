package main

import (
	"context"
	"net"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

const greetingMessage = "Hello. Speak, and I will listen to you."
const partingMessage = "Good bye. Thanks for calling."
const timeoutMessage = "Sorry, your time is up."

var keyPhrases = []string{
	"bye",
	"goodbye",
	"hangup",
}

type stateFn func(ctx context.Context) (stateFn, error)

// App is the base state machine application for the call
type App struct {
	c  net.Conn
	id uuid.UUID
}

func (a *App) Run(ctx context.Context) (err error) {
	for next := a.rootMenu; next != nil; {
		if ctx.Err() != nil {
			return nil
		}
		if next, err = next(ctx); err != nil {
			return err
		}
	}
	return ErrHangup
}

func (a *App) rootMenu(ctx context.Context) (stateFn, error) {
	if err := speak(ctx, a.c, greetingMessage); err != nil {
		return nil, errors.Wrap(err, "failed to send greeting to asterisk")
	}

	cmd, err := recognizeRequest(ctx, a.c)
	if err != nil {
		return a.listenFailure, nil
	}

	if containsAny(cmd, "time") {
		return a.tellTime, nil
	}
	if containsAny(cmd, "laugh", "joke") {
		return a.tellJoke, nil
	}
	if containsAny(cmd, "bye", "hangup", "hang up") {
		return nil, ErrHangup
	}

	return a.rootMenu, nil
}

func (a *App) listenFailure(ctx context.Context) (stateFn, error) {
	err := speak(ctx, a.c, "Sorry, I failed to listen")
	return a.rootMenu, err
}
