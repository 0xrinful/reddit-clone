package server

import (
	"net/http"

	"github.com/0xrinful/rush"

	"github.com/0xrinful/reddit-clone/internal/middleware"
	"github.com/0xrinful/reddit-clone/internal/shared/response"
)

func (s *Server) setupRoutes(responder *response.Responder) http.Handler {
	r := rush.New()
	m := middleware.New(responder, s.cfg)
	r.Use(m.Recover)
	r.Use(m.Authenticate(s.tokensSvc))

	r.NotFound = http.HandlerFunc(responder.NotFound)
	r.MethodNotAllowed = http.HandlerFunc(responder.MethodNotAllowed)

	r.Route("/api/v1", func(r *rush.Router) {
		// ─── AUTH ─────────────────────────────────────────────────────────
		r.Route("/auth", func(r *rush.Router) {
			r.Group(func(r *rush.Router) {
				r.Use(m.RequireUnauth)
				r.With(m.RegisterLimit()).Post("/register", s.auth.RegisterUser)
				r.With(m.AuthLimit()).Post("/login", s.auth.Login)
			})

			r.Group(func(r *rush.Router) {
				r.Use(m.AuthLimit())
				r.Post("/refresh", s.auth.Refresh)
				r.Post("/logout", s.auth.Logout)
			})

			r.Route("/email", func(r *rush.Router) {
				r.Use(m.StrictLimit())
				r.Post("/email/verify", s.auth.ActivateUser)
				r.Post("/verify/resend", s.auth.SendActivationEmail)
			})
		})

		// ─── COMMUNITIES ──────────────────────────────────────────────────
		r.Route("/communities", func(r *rush.Router) {
			r.Get("/", nil)
			r.Post("/", nil)

			r.Route("/{community_name}", func(r *rush.Router) {
				r.Use(m.LoadCommunity(s.communitiesSvc))

				r.With(m.ReadLimit()).Get("/posts", s.posts.List)
				r.With(m.WriteLimit(), m.RequireAuth).Post("/posts", s.posts.Create)
			})
		})

		// ─── POSTS ────────────────────────────────────────────────────────
		r.Route("/posts", func(r *rush.Router) {
			r.With(m.ReadLimit()).Get("/{id}", s.posts.Get)
			r.Group(func(r *rush.Router) {
				r.Use(m.WriteLimit(), m.RequireAuth)
				r.Delete("/{id}", s.posts.Delete)
				r.Patch("/{id}", s.posts.Update)
			})
		})
	})

	return r
}
