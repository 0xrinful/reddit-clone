package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/0xrinful/reddit-clone/internal/auth"
	"github.com/0xrinful/reddit-clone/internal/bans"
	"github.com/0xrinful/reddit-clone/internal/communities"
	"github.com/0xrinful/reddit-clone/internal/config"
	"github.com/0xrinful/reddit-clone/internal/members"
	"github.com/0xrinful/reddit-clone/internal/posts"
	"github.com/0xrinful/reddit-clone/internal/shared/background"
	"github.com/0xrinful/reddit-clone/internal/shared/response"
	"github.com/0xrinful/reddit-clone/internal/tokens"
	"github.com/0xrinful/reddit-clone/internal/users"
)

type Server struct {
	cfg        config.Config
	logger     *slog.Logger
	httpServer *http.Server
	bg         *background.Worker

	// services (for middleware)
	communitiesSvc communities.Service
	tokensSvc      tokens.Service

	// handlers (for routes)
	posts       *posts.Handler
	communities *communities.Handler
	users       *users.Handler
	auth        *auth.Handler
	members     *members.Handler
	bans        *bans.Handler
}

func New(
	cfg config.Config,
	logger *slog.Logger,
	bg *background.Worker,
	communitiesSvc communities.Service,
	postsSvc posts.Service,
	usersSvc users.Service,
	authSvc auth.Service,
	tokensSvc tokens.Service,
	membersSvc members.Service,
	bansSvc bans.Service,
) *Server {
	responder := response.NewResponder(logger)

	postsHandler := posts.NewHandler(postsSvc, responder)
	communitiesHandler := communities.NewHandler(communitiesSvc, responder)
	usersHandler := users.NewHandler(usersSvc, responder)
	authHandler := auth.NewHandler(authSvc, tokensSvc, responder)
	membersHandler := members.NewHandler(membersSvc, responder)
	bansHandler := bans.NewHandler(bansSvc, responder)

	server := &Server{
		cfg:    cfg,
		logger: logger,
		bg:     bg,

		communitiesSvc: communitiesSvc,
		tokensSvc:      tokensSvc,

		posts:       postsHandler,
		communities: communitiesHandler,
		users:       usersHandler,
		auth:        authHandler,
		members:     membersHandler,
		bans:        bansHandler,
	}

	router := server.setupRoutes(responder)

	// bridge slog → *log.Logger for http.Server
	errLog := slog.NewLogLogger(logger.Handler(), slog.LevelError)

	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ErrorLog:     errLog,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server
}

func (s *Server) Start() error {
	s.logger.Info("server starting", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	s.logger.Info("waiting for background tasks to finish...")
	s.bg.Wait()
	s.logger.Info("all tasks finished, exiting")
	return nil
}
