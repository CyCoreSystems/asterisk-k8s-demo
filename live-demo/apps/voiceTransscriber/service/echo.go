package main

import (
	"context"

	"github.com/pkg/errors"
)

func (a *App) echoStart(ctx context.Context) (stateFn, error) {
	if err := speak(ctx, a.c, "go ahead.  say cancel or menu to exit"); err != nil {
		return nil, errors.Wrap(err, "failed to send message to asterisk")
	}
	return a.echo, nil
}

func (a *App) echo(ctx context.Context) (stateFn, error) {
	cmd, err := recognizeRequest(ctx, a.c)
	if err != nil {
		return a.listenFailure, nil
	}

	if containsAny(cmd, "bye", "hangup", "hang up") {
		return nil, ErrHangup
	}
	if containsAny(cmd, "cancel", "menu") {
		return a.rootMenu, nil
	}
	err = speak(ctx, a.c, cmd)
	return a.echo, err
}
