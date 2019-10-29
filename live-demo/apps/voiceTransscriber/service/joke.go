package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

var jokes = []string{
	`I asked God for a bike, but I know God doesn't work that way so I stole a bike and asked for forgiveness.`,
	`I hate Russian dolls, they're so full of themselves.`,
	`Throwing acid is wrong, in some people's eyes.`,
	`The first time I got a universal remote control I thought to myself, "This changes everything".`,
	`Say what you want about deaf people...`,
	`I've spent the last four years looking for my ex-girlfriend's killer, but no-one will do it.`,
	`I refused to believe my road worker father was stealing from his job, but when I got home all the signs were there.`,
	`I recently decided to sell my vacuum cleaner as all it was doing was gathering dust.`,
	`You can never lose a homing pigeon - if your homing pigeon doesn't come back what you've lost is a pigeon.`,
	`My girlfriend told me to go out and get something that makes her look sexy... so I got drunk.`,
	`Don't you hate it when someone answers their own questions? I do.`,
	`As I watched the dog chasing his tail I thought "Dogs are easily amused", then I realized I was watching the dog chasing his tail.`,
	`PMS jokes are not funny or appropriate. Period!`,
	`Gambling addiction hotlines would do so much better if every fifth caller was a winner.`,
	`Where there's a will, there's a relative.`,
	`Hedgehogs, eh? Why can't they just share the hedge?`,
	`Just because nobody complains doesn't mean all parachutes are perfect.`,
	`To the man on crutches, dressed in camouflage, who stole my wallet - you can hide, but you can't run.`,
	`Velcro - what a rip-off!`,
	`My friend keeps trying to convince me that he's a compulsive liar but I don't believe him.`,
	`It’s always hard to explain puns to kleptomaniacs because they’re always taking things literally.`,
}

func init() {
	rand.Seed(time.Now().Unix())
}

func (a *App) tellJoke(ctx context.Context) (stateFn, error) {
	if err := speak(ctx, a.c, jokes[rand.Intn(len(jokes))]); err != nil {
		return nil, errors.Wrap(err, "failed to send message to asterisk")
	}
	return a.rootMenu, nil
}
