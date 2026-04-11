package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	lastSeen time.Time
}

func (m *Middleware) startClientCleanup() {
	if !m.config.Limiter.Enabled {
		return
	}
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		for k, c := range m.clients {
			if time.Since(c.lastSeen) > 5*time.Minute {
				delete(m.clients, k)
			}
		}
		m.mu.Unlock()
	}
}

func (m *Middleware) RateLimit(
	name string,
	requests int,
	per time.Duration,
) func(http.Handler) http.Handler {
	if !m.config.Limiter.Enabled {
		return func(next http.Handler) http.Handler { return next }
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			k := key(r)

			m.mu.RLock()
			c, exists := m.clients[k]
			m.mu.RUnlock()

			if !exists {
				m.mu.Lock()
				if c, exists = m.clients[k]; !exists {
					c = &client{limiters: make(map[string]*rate.Limiter), lastSeen: time.Now()}
					m.clients[k] = c
				}
				m.mu.Unlock()
			}

			c.mu.Lock()
			c.lastSeen = time.Now()
			limiter, exists := c.limiters[name]
			if !exists {
				limiter = rate.NewLimiter(rate.Limit(requests)/rate.Limit(per.Seconds()), requests)
				c.limiters[name] = limiter
			}
			c.mu.Unlock()

			if !limiter.Allow() {
				m.responder.TooManyRequests(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (m *Middleware) ReadLimit() func(http.Handler) http.Handler {
	return m.RateLimit("read", 200, time.Minute)
}

func (m *Middleware) WriteLimit() func(http.Handler) http.Handler {
	return m.RateLimit("write", 30, time.Minute)
}

func (m *Middleware) StrictLimit() func(http.Handler) http.Handler {
	return m.RateLimit("strict", 5, time.Minute*10)
}

func key(r *http.Request) string {
	// TODO: add user session later if authenticated instead of ip
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
