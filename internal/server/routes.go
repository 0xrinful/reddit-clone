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

	r.NotFound = http.HandlerFunc(responder.NotFound)
	r.MethodNotAllowed = http.HandlerFunc(responder.MethodNotAllowed)

	r.Route("/api/v1", func(r *rush.Router) {
		r.Route("/r/{community_name}", func(r *rush.Router) {
			r.Use(middleware.LoadCommunity(communitySvc))
			postsHanlder.RegisterRoutes(r)
		})
	})

	return r
}
