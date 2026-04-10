package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/0xrinful/reddit-clone/internal/communities"
	"github.com/0xrinful/reddit-clone/internal/config"
	"github.com/0xrinful/reddit-clone/internal/middleware"
	"github.com/0xrinful/reddit-clone/internal/posts"
	"github.com/0xrinful/reddit-clone/internal/shared/background"
	"github.com/0xrinful/reddit-clone/internal/shared/response"
	"github.com/0xrinful/reddit-clone/internal/users"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	background *background.Worker
}

func New(
	cfg config.Config,
	logger *slog.Logger,
	background *background.Worker,
	communitiesSvc communities.Service,
	postsSvc posts.Service,
	usersSvc users.Service,
) *Server {
	responder := response.NewResponder(logger)
	middleware := middleware.New(responder, cfg)

	postsHandler := posts.NewHandler(postsSvc, responder)
	usersHandler := users.NewHandler(usersSvc, responder)
	router := setupRoutes(responder, middleware, communitiesSvc, postsHandler, usersHandler)

	// bridge slog → *log.Logger for http.Server
	errLog := slog.NewLogLogger(logger.Handler(), slog.LevelError)

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      router,
			ErrorLog:     errLog,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		logger:     logger,
		background: background,
	}
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
	s.background.Wait()
	s.logger.Info("all tasks finished, exiting")
	return nil
}
