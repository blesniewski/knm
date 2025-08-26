package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/blesniewski/knm/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	err := app.Run(ctx)
	if err != nil {
		panic(err)
	}
}
