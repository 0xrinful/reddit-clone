package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func (m *Middleware) RateLimit(requests int, per time.Duration) func(http.Handler) http.Handler {
	if !m.config.Limiter.Enabled {
		return func(next http.Handler) http.Handler { return next }
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()
			for key, client := range clients {
				if time.Since(client.lastSeen) > time.Minute*5 {
					delete(clients, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			k := key(r)

			mu.Lock()
			if _, found := clients[k]; !found {
				clients[k] = &client{
					limiter: rate.NewLimiter(rate.Every(per/time.Duration(requests)), requests),
				}
			}
			c := clients[k]
			c.lastSeen = time.Now()

			if !c.limiter.Allow() {
				mu.Unlock()
				m.responder.TooManyRequests(w, r)
				return
			}
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func (m *Middleware) ReadLimit() func(http.Handler) http.Handler {
	return m.RateLimit(200, time.Minute)
}

func (m *Middleware) WriteLimit() func(http.Handler) http.Handler {
	return m.RateLimit(30, time.Minute)
}

func (m *Middleware) AuthLimit() func(http.Handler) http.Handler {
	return m.RateLimit(5, time.Minute)
}

func key(r *http.Request) string {
	// TODO: add user session later if authenticated instead of ip
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
