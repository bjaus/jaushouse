package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env             string        `envconfig:"ENV" default:"local"`
	Port            int           `envconfig:"PORT" default:"8080"`
	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"60s"`
}

func main() {
	ctx := context.Background()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if envfile := os.Getenv("ENV_FILE"); envfile == "" {
		if err := godotenv.Load(envfile); err != nil {
			log.Warn("Failed to load .env file", "error", err)
		}
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	if err := run(ctx, log, cfg); err != nil {
		log.Error("Error running service", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *slog.Logger, cfg Config) error {
	log.Info("Service started")

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := http.Server{
		Handler:      r,
		Addr:         ":" + strconv.Itoa(cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go gracefulShutdown(ctx, &srv, cfg.ShutdownTimeout, log)

	log.Info("Starting server", "address", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	log.Info("Service stopped")

	return nil
}

func gracefulShutdown(ctx context.Context, srv *http.Server, timeout time.Duration, log *slog.Logger) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Error("Shutdown timeout exceeded", "error", err)
		}
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Failed to shutdown server", "error", err)
	}
}
