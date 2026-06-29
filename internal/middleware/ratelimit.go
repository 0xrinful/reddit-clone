package middleware

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/0xrinful/reddit-clone/internal/shared/request"
)

type limiterEntry struct {
	limiter   *rate.Limiter
	expiresAt time.Time
}

func newEntry(requests int, per time.Duration, ttl time.Duration) *limiterEntry {
	return &limiterEntry{
		limiter:   rate.NewLimiter(rate.Limit(requests)/rate.Limit(per.Seconds()), requests),
		expiresAt: time.Now().Add(ttl),
	}
}

// TODO: add number of blocked users limit for the ip [block ip which has too many bad users]
type client struct {
	mu       sync.Mutex
	limiters map[string]*limiterEntry
}

func (m *Middleware) startClientCleanup() {
	if !m.config.Limiter.Enabled {
		return
	}
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		for k, c := range m.clients {
			// Remove expired rate limit entries for this client
			c.mu.Lock()
			for name, entry := range c.limiters {
				if time.Now().After(entry.expiresAt) {
					delete(c.limiters, name)
				}
			}
			c.mu.Unlock()

			// Remove client if it has no active rate limit entries
			if len(c.limiters) == 0 {
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
	ttl time.Duration,
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
					c = &client{limiters: make(map[string]*limiterEntry)}
					m.clients[k] = c
				}
				m.mu.Unlock()
			}

			c.mu.Lock()
			entry, exists := c.limiters[name]
			if !exists || time.Now().After(entry.expiresAt) {
				entry = newEntry(requests, per, ttl)
				c.limiters[name] = entry
			}
			allowed := entry.limiter.Allow()
			if allowed {
				entry.expiresAt = time.Now().Add(ttl)
			}
			c.mu.Unlock()

			if !allowed {
				m.responder.TooManyRequests(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (m *Middleware) ReadLimit() func(http.Handler) http.Handler {
	return m.RateLimit("read", 120, time.Minute, 5*time.Minute)
}

func (m *Middleware) WriteLimit() func(http.Handler) http.Handler {
	return m.RateLimit("write", 30, time.Minute, 5*time.Minute)
}

func (m *Middleware) AuthLimit() func(http.Handler) http.Handler {
	return m.RateLimit("auth", 10, 15*time.Minute, 15*time.Minute)
}

func (m *Middleware) StrictLimit() func(http.Handler) http.Handler {
	return m.RateLimit("strict", 3, time.Hour*1, 1*time.Hour)
}

func (m *Middleware) RegisterLimit() func(http.Handler) http.Handler {
	return m.RateLimit("register", 5, time.Hour*12, 24*time.Hour)
}

// TODO: support X-Forwarded-For, X-Real-IP, ...
func key(r *http.Request) string {
	user, authenticated := request.GetUser(r)
	if authenticated {
		return "user:" + strconv.FormatInt(user.ID, 10)
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return "ip:" + ip
}
