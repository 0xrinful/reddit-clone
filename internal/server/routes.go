package server

import (
	"net/http"

	"github.com/0xrinful/rush"

	"github.com/0xrinful/reddit-clone/internal/communities"
	"github.com/0xrinful/reddit-clone/internal/middleware"
	"github.com/0xrinful/reddit-clone/internal/posts"
	"github.com/0xrinful/reddit-clone/internal/shared/response"
)

func setupRoutes(
	responder *response.Responder,
	middleware *middleware.Middleware,
	communitySvc communities.Service,
	postsHanlder *posts.Handler,
) http.Handler {
	r := rush.New()
	r.Use(middleware.Recover)

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
				r.Use(middleware.WriteLimit())
				r.Post("/posts", postsHanlder.Create)
				r.Delete("/posts/{id}", postsHanlder.Delete)
				r.Patch("/posts/{id}", postsHanlder.Update)
			})
		})
	})

	return r
}
