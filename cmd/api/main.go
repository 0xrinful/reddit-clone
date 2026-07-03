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

	"github.com/joho/godotenv"

	"github.com/0xrinful/reddit-clone/internal/auth"
	"github.com/0xrinful/reddit-clone/internal/communities"
	"github.com/0xrinful/reddit-clone/internal/config"
	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/members"
	"github.com/0xrinful/reddit-clone/internal/posts"
	"github.com/0xrinful/reddit-clone/internal/server"
	"github.com/0xrinful/reddit-clone/internal/shared/background"
	"github.com/0xrinful/reddit-clone/internal/shared/mailer"
	"github.com/0xrinful/reddit-clone/internal/tokens"
	"github.com/0xrinful/reddit-clone/internal/users"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Warn("no .env file loaded", "err", err)
	}

	cfg := config.Load()
	mailer := mailer.New(cfg.SMTP)
	bg := background.New(logger)

	db, err := database.Open(cfg.DB)
	if err != nil {
		logger.Error("db connection failed", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	communitiesRepo := communities.NewRepository(db)
	postsRepo := posts.NewRepository(db)
	usersRepo := users.NewRepository(db)
	tokensRepo := tokens.NewRepository(db)
	membersRepo := members.NewRepository(db)

	communitiesSvc := communities.NewService(db, communitiesRepo, membersRepo)
	postsSvc := posts.NewService(postsRepo)
	usersSvc := users.NewService(usersRepo)
	authSvc := auth.NewService(db, usersRepo, tokensRepo, mailer, logger, bg)
	tokensSvc := tokens.NewService(db, tokensRepo, cfg.JWT)

	srv := server.New(cfg, logger, bg, communitiesSvc, postsSvc, usersSvc, authSvc, tokensSvc)

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
