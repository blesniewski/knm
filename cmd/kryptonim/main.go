package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/blesniewski/knm/internal/api"
	"github.com/blesniewski/knm/internal/clients/cryptoexchange"
	"github.com/blesniewski/knm/internal/clients/oxr"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	oxrClient, err := oxr.NewClient(ctx, cfg.Orx.BaseURL, cfg.Orx.AppID)
	if err != nil {
		panic(err)
	}
	cryptoClient := cryptoexchange.NewClient()
	httpServer := api.NewServer(oxrClient, cryptoClient)
	httpServer.Run(cfg.ListenAddr)
}
