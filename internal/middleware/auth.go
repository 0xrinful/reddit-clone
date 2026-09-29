package middleware

import (
	"net/http"
	"strings"

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
				m.responder.HandleServiceError(w, err)
				return
			}

			r = request.WithUser(r, &request.UserCtx{ID: claims.UserID})
			next.ServeHTTP(w, r)
		})
	}
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, authenticated := request.GetUser(r)
		if !authenticated {
			m.responder.Unauthorized(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireUnauth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, authenticated := request.GetUser(r)
		if authenticated {
			m.responder.Forbidden(w, "already_authenticated", "you are already authenticated")
			return
		}
		next.ServeHTTP(w, r)
	})
}
