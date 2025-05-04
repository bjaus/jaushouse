package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"pkg/config"
	"pkg/web"
)

func setupHandler() http.Handler {
	r := web.DefaultRouter()
	return r
}

type Config struct {
	Env             string        `envconfig:"ENV" default:"local"`
	Port            int           `envconfig:"PORT" default:"8080"`
	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"60s"`
}

func run(ctx context.Context, log *slog.Logger, cfg Config) error {
	log.Info("Service started")

	srv := http.Server{
		Handler:      setupHandler(),
		Addr:         ":" + strconv.Itoa(cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go web.GracefulShutdown(ctx, &srv, cfg.ShutdownTimeout)

	log.Info("Starting server", "address", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	log.Info("Service stopped")

	return nil
}

func main() {
	ctx := context.Background()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load[Config]()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	if err := run(ctx, log, cfg); err != nil {
		log.Error("Error running service", "error", err)
		os.Exit(1)
	}
}
