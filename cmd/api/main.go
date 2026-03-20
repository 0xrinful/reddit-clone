package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xrinful/reddit-clone/internal/config"
	"github.com/0xrinful/reddit-clone/internal/posts"
	"github.com/0xrinful/reddit-clone/internal/server"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	// repos
	postsRepo := posts.NewRepository(nil)

	// services
	postsSvc := posts.NewService(postsRepo)

	// server
	srv := server.New(cfg, postsSvc, logger)

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	select {
	case err := <-errCh:
		// server crashed unexpectedly
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "err", err)
		}

	case <-ctx.Done():
		logger.Info("shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown error", "err", err)
		}

		// wait for ListenAndServe to return
		err := <-errCh
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "err", err)
		}
	}

	logger.Info("server stopped")
}
