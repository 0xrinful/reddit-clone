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

		rw := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		method := r.Method
		path := r.URL.Path
		status := rw.status
		duration := time.Since(start).Round(time.Microsecond)

		if m.config.Logging.IsProduction {
			logRequestStructured(method, path, status, duration)
		} else {
			logRequestColored(method, path, status, duration)
		}
	})
}

func logRequestStructured(method, path string, status int, duration time.Duration) {
	slog.Info("request", "method", method, "path", path, "status", status, "duration", duration)
}

func logRequestColored(method, path string, status int, duration time.Duration) {
	slog.Info(
		fmt.Sprintf(
			"request: \033[33m%s\033[0m \033[36m%s\033[0m \033[35m%d\033[0m \033[32m%s\033[0m",
			method, path, status, duration,
		),
	)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (w *statusRecorder) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusRecorder) Unwrap() http.ResponseWriter { return w.ResponseWriter }
