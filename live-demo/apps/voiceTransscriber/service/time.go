package main

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

func (a *App) tellTime(ctx context.Context) (stateFn, error) {
	if err := speak(ctx, a.c, time.Now().Format("15 04")); err != nil {
		return nil, errors.Wrap(err, "failed to send message to asterisk")
	}
	return a.rootMenu, nil
}
