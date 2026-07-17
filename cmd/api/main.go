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

	"github.com/0xrinful/reddit-clone/internal/auth"
	"github.com/0xrinful/reddit-clone/internal/authorization"
	"github.com/0xrinful/reddit-clone/internal/bans"
	"github.com/0xrinful/reddit-clone/internal/communities"
	"github.com/0xrinful/reddit-clone/internal/config"
	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/members"
	"github.com/0xrinful/reddit-clone/internal/posts"
	"github.com/0xrinful/reddit-clone/internal/server"
	"github.com/0xrinful/reddit-clone/internal/shared/background"
	"github.com/0xrinful/reddit-clone/internal/shared/logger"
	"github.com/0xrinful/reddit-clone/internal/shared/mailer"
	"github.com/0xrinful/reddit-clone/internal/tokens"
	"github.com/0xrinful/reddit-clone/internal/users"
	"github.com/0xrinful/reddit-clone/internal/votes"
)

func main() {
	cfg := config.Load()

	logger := logger.New(cfg.Logging)
	slog.SetDefault(logger)

	mailer := mailer.New(cfg.SMTP)
	bg := background.New()

	db, err := database.Open(cfg.DB)
	if err != nil {
		slog.Error("db connection failed", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("database connection pool established")

	communitiesRepo := communities.NewRepository(db)
	postsRepo := posts.NewRepository(db)
	usersRepo := users.NewRepository(db)
	tokensRepo := tokens.NewRepository(db)
	membersRepo := members.NewRepository(db)
	bansRepo := bans.NewRepository(db)
	votesRepo := votes.NewRepository(db)

	authzSvc := authorization.NewService(bansRepo, membersRepo)
	communitiesSvc := communities.NewService(db, authzSvc, communitiesRepo, membersRepo)
	postsSvc := posts.NewService(db, authzSvc, postsRepo, votesRepo)
	usersSvc := users.NewService(usersRepo)
	authSvc := auth.NewService(db, usersRepo, tokensRepo, mailer, bg)
	tokensSvc := tokens.NewService(db, tokensRepo, cfg.JWT)
	membersSvc := members.NewService(authzSvc, membersRepo, usersRepo)
	bansSvc := bans.NewService(db, authzSvc, bansRepo, usersRepo, postsRepo)
	votesSvc := votes.NewService(db, votesRepo, postsRepo)

	srv := server.New(
		cfg, bg,
		communitiesSvc,
		postsSvc, usersSvc, authSvc,
		tokensSvc, membersSvc, bansSvc,
		votesSvc,
	)

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
			slog.Error("server error", "err", err)
		}

	case <-ctx.Done():
		slog.Info("shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutdown error", "err", err)
		}

		// wait for ListenAndServe to return
		err := <-errCh
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
		}
	}

	slog.Info("server stopped")
}
