package middleware

import (
	"log"
	"net/http"
	"time"
)

func (m *Middleware) Logger(next http.Handler) http.Handler {
	if !m.config.Logging.RequestLogging {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		log.Printf(
			"\033[33m%s\033[0m \033[36m%s\033[0m \033[32m%s\033[0m",
			r.Method, r.URL.Path, time.Since(start),
		)
	})
}
