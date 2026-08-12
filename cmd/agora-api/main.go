// Command agora-api is the ÁGORA election API server.
//
// At this stage it is the service skeleton and nothing more: configuration,
// structured logging, a health endpoint, timeouts and graceful shutdown. The
// election, ballot submission and tally endpoints are the work of the MVP phase
// and belong on top of this, not inside it.
//
// The pieces that are here rather than deferred are the ones that are painful to
// retrofit: a server without ReadHeaderTimeout is a slowloris target, and a
// server without graceful shutdown drops in-flight requests on every deploy.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/felicabrera/agora/internal/config"
	"github.com/felicabrera/agora/internal/version"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "agora-api: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	logger := newLogger(cfg.LogLevel)
	build := version.Read()
	logger.Info("starting agora-api", "addr", cfg.Addr, "build", build.String())

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           routes(logger),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	// Shut down on SIGINT/SIGTERM, giving in-flight requests a bounded window to
	// finish rather than cutting them off.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("listening: %w", err)
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		logger.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutting down: %w", err)
		}
		return nil
	}
}

func routes(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if _, err := fmt.Fprintf(w, `{"status":"ok","build":%q}`+"\n", version.Read().String()); err != nil {
			logger.Error("writing health response", "error", err)
		}
	})
	return mux
}

func newLogger(level string) *slog.Logger {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		l = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}
