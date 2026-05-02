package server

import (
	"net/http"

	"github.com/0xrinful/rush"

	"github.com/0xrinful/reddit-clone/internal/auth"
	"github.com/0xrinful/reddit-clone/internal/communities"
	"github.com/0xrinful/reddit-clone/internal/middleware"
	"github.com/0xrinful/reddit-clone/internal/posts"
	"github.com/0xrinful/reddit-clone/internal/shared/response"
	"github.com/0xrinful/reddit-clone/internal/tokens"
	"github.com/0xrinful/reddit-clone/internal/users"
)

func setupRoutes(
	responder *response.Responder,
	middleware *middleware.Middleware,
	communitySvc communities.Service,
	tokensSvc tokens.Service,
	postsHanlder *posts.Handler,
	usersHanlder *users.Handler,
	authHandler *auth.Handler,
) http.Handler {
	r := rush.New()
	r.Use(middleware.Recover)
	r.Use(middleware.Authenticate(tokensSvc))

	r.NotFound = http.HandlerFunc(responder.NotFound)
	r.MethodNotAllowed = http.HandlerFunc(responder.MethodNotAllowed)

	r.Route("/api/v1", func(r *rush.Router) {
		r.Route("/r/{community_name}", func(r *rush.Router) {
			r.Use(middleware.LoadCommunity(communitySvc))

			r.Group(func(r *rush.Router) {
				r.Use(middleware.ReadLimit())
				r.Get("/posts", postsHanlder.List)
				r.Get("/posts/{id}", postsHanlder.Get)
			})

			r.Group(func(r *rush.Router) {
				r.Use(middleware.RequireAuth)
				r.Use(middleware.WriteLimit())
				r.Post("/posts", postsHanlder.Create)
				r.Delete("/posts/{id}", postsHanlder.Delete)
				r.Patch("/posts/{id}", postsHanlder.Update)
			})
		})

		r.Route("/auth", func(r *rush.Router) {
			r.Group(func(r *rush.Router) {
				r.Use(middleware.RequireUnauth)
				r.With(middleware.AuthLimit()).Post("/login", authHandler.Login)
				r.With(middleware.RegisterLimit()).Post("/register", authHandler.RegisterUser)
			})

			r.Group(func(r *rush.Router) {
				r.Use(middleware.AuthLimit())
				r.Post("/refresh", authHandler.Refresh)
				r.Post("/logout", authHandler.Logout)
				r.Post("/email/verify", authHandler.ActivateUser)
				r.Post("/email/verify/resend", authHandler.SendActivationEmail)
			})
		})
	})

	return r
}
