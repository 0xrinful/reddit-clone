package middleware

import (
	"fmt"
	"log/slog"
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

		method := r.Method
		path := r.URL.Path
		duration := time.Since(start).Round(time.Microsecond)

		if m.config.Logging.IsProduction {
			slog.Info("request", "method", method, "path", path, "duration", duration)
		} else {
			slog.Info("request: " +
				fmt.Sprintf(
					"\033[33m%s\033[0m \033[36m%s\033[0m \033[32m%s\033[0m",
					method, path, duration,
				),
			)
		}
	})
}
