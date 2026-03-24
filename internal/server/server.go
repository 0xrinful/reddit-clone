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
	"github.com/0xrinful/reddit-clone/internal/shared/response"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(
	cfg config.Config,
	communitiesSvc communities.Service,
	postsSvc posts.Service,
	logger *slog.Logger,
) *Server {
	responder := response.NewResponder(logger)
	middleware := middleware.New(responder)

	postsHandler := posts.NewHandler(postsSvc, responder)
	router := setupRoutes(responder, middleware, communitiesSvc, postsHandler)

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
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("server starting", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
