package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/alexPavlikov/go-wiki/internal/config"
	router "github.com/alexPavlikov/go-wiki/internal/server"
	"github.com/alexPavlikov/go-wiki/internal/server/locations"
	"github.com/alexPavlikov/go-wiki/internal/server/service"
)

func Run() error {

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed load config", "error", err)
		return fmt.Errorf("failed load config: %w", err)
	}

	slog.Info("starting application on", "path", cfg.Server.ToString())

	handler := locations.New(service.Service{})

	rtr := router.New(handler)

	srv := rtr.Build()

	// load http server
	if err := http.ListenAndServe(cfg.Server.ToString(), srv); err != nil {
		slog.Error("listen and serve server error", "error", err.Error())
		return err
	}

	return nil
}
