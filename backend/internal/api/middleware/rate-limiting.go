package middleware

import (
	response "dukkanim-api/internal/platform/http_response"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

/*
	ratelimiter.Config {
		RequestPerTimeFrame int
		TimeFrame int -> time.Second * 5 e.g.
		Enabled bool
	}
*/

type RateLimitingMiddleware struct {
	counter safeCounter
	RatelimiterConfig
	windowStart time.Time
}

type RatelimiterConfig struct {
	RequestPerTimeFrame int
	TimeFrameSeconds    int
}

type safeCounter struct {
	mu      sync.Mutex
	counter map[string]int
}

func NewRateLimitingMiddleware(cfg RatelimiterConfig) *RateLimitingMiddleware {

	rl := &RateLimitingMiddleware{
		counter: safeCounter{
			counter: map[string]int{},
		},
		RatelimiterConfig: cfg,
		windowStart:       time.Now(),
	}

	go func() {
		clock := time.NewTicker(time.Second * time.Duration(cfg.TimeFrameSeconds))
		for range clock.C {

			rl.counter.reset()
			rl.windowStart = time.Now()

		}
	}()

	return rl
}

func (r *safeCounter) inc(ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counter[ip]++
}
func (r *safeCounter) reset(ip ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(ip) == 0 {
		r.counter = map[string]int{}
		return
	}
	r.counter[ip[0]] = 0
}
func (r *safeCounter) value(ip string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.counter[ip]
}

func (m *RateLimitingMiddleware) RateLimiting(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// use configurable limit

		// Extract client IP from the request
		clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Server Failed to Read Client IP", http.StatusInternalServerError)
			return
		}

		// check current request count for that IP in current time window
		m.counter.inc(clientIP)
		count := m.counter.value(clientIP)

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d requests/per minute", m.RequestPerTimeFrame))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, m.RequestPerTimeFrame-count)))

		if count > m.RequestPerTimeFrame {
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", m.windowStart.Add(time.Duration(m.TimeFrameSeconds)*time.Second).Unix()))
			response.RespondWithProblemDetails(w, r.Context(), http.StatusTooManyRequests, "Request Rate Limited", "RATE_LIMIT", nil)
			return
		}

		// if under limit increment counter than next.ServeHTTP
		// increment counter

		// Set
		// X-RateLimit-Limit
		// & X-RateLimit-Remanining
		// and X-RateLimit-Reset Headers

		next.ServeHTTP(w, r)

		// if over limit return http 429 with X-RateLimit-Retry-After

	})
}
