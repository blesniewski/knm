package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/blesniewski/knm/internal/clients/cryptoexchange"
	"github.com/blesniewski/knm/internal/clients/oxr"
	api "github.com/blesniewski/knm/internal/server"
)

func Run(ctx context.Context) error {
	cfg, err := NewConfig()
	if err != nil {
		return err
	}

	oxrClient, err := oxr.NewClient(ctx, cfg.Orx.BaseURL, cfg.Orx.AppID)
	if err != nil {
		return fmt.Errorf("failed to create oxr client: %w", err)
	}

	cryptoClient := cryptoexchange.NewClient()
	httpServer := api.New(oxrClient, cryptoClient)
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- httpServer.Run(cfg.ListenAddr)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
		return nil
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to run http server: %w", err)
		}
		return nil
	}
}
