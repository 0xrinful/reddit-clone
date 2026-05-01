package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/0xrinful/reddit-clone/internal/shared/errs"
	"github.com/0xrinful/reddit-clone/internal/shared/request"
	"github.com/0xrinful/reddit-clone/internal/tokens"
)

func (m *Middleware) Authenticate(svc tokens.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				m.responder.InvalidToken(w)
				return
			}

			claims, err := svc.VerifyAccessToken(parts[1])
			if err != nil {
				switch {
				case errors.Is(err, errs.ErrInvalidToken):
					m.responder.InvalidToken(w)
				default:
					m.responder.ServerError(w, err)
				}
				return
			}

			u := &request.UserCtx{ID: claims.UserID}
			r = request.WithUser(r, u)
			next.ServeHTTP(w, r)
		})
	}
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := request.GetUser(r)
		if !ok {
			m.responder.Unauthorized(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
