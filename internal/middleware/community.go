package middleware

import (
	"errors"
	"net/http"

	"github.com/0xrinful/reddit-clone/internal/communities"
	"github.com/0xrinful/reddit-clone/internal/shared/apperr"
	"github.com/0xrinful/reddit-clone/internal/shared/request"
)

func (m *Middleware) LoadCommunity(svc communities.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			name := r.PathValue("community_name")
			c, err := svc.GetByName(r.Context(), name)
			if err != nil {
				switch {
				case errors.Is(err, apperr.ErrNotFound):
					m.responder.NotFound(w, r)
				default:
					m.responder.ServerError(w, err)
				}
				return
			}
			r = request.WithCommunity(r, &request.CommunityCtx{ID: c.ID, Name: c.Name})
			next.ServeHTTP(w, r)
		})
	}
}
