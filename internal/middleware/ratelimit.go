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

	var clients sync.Map

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			clients.Range(func(key, value any) bool {
				c := value.(*client)
				if time.Since(c.lastSeen) > time.Minute*5 {
					clients.Delete(key)
				}
				return true
			})
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			k := key(r)

			val, loaded := clients.Load(k)
			if !loaded {
				newClient := &client{
					limiter: rate.NewLimiter(
						rate.Limit(requests)/rate.Limit(per.Seconds()),
						requests,
					),
					lastSeen: time.Now(),
				}

				val, _ = clients.LoadOrStore(k, newClient)
			}
			c := val.(*client)
			c.lastSeen = time.Now()

			if !c.limiter.Allow() {
				m.responder.TooManyRequests(w, r)
				return
			}

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
