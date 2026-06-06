package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/masterfabric-go/masterfabric/internal/shared/response"
)

// RateLimiter is a simple in-memory fixed-window limiter keyed by client IP.
// It throttles sensitive endpoints (login/register) to slow brute-force and
// abuse. For multi-instance deployments front this with a shared limiter
// (Redis/Nginx/WAF); this provides per-instance defense in depth (KVKK Art. 12
// requires "appropriate technical measures" against unlawful access).
type RateLimiter struct {
	mu       sync.Mutex
	counters map[string]*rlWindow
	limit    int
	window   time.Duration
}

type rlWindow struct {
	count int
	reset time.Time
}

// NewRateLimiter builds a limiter allowing `limit` requests per `per` window.
func NewRateLimiter(limit int, per time.Duration) *RateLimiter {
	if limit <= 0 {
		limit = 10
	}
	if per <= 0 {
		per = time.Minute
	}
	rl := &RateLimiter{
		counters: make(map[string]*rlWindow),
		limit:    limit,
		window:   per,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	t := time.NewTicker(rl.window)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		rl.mu.Lock()
		for k, w := range rl.counters {
			if now.After(w.reset) {
				delete(rl.counters, k)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) allow(key string) (bool, time.Duration) {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()
	w, ok := rl.counters[key]
	if !ok || now.After(w.reset) {
		rl.counters[key] = &rlWindow{count: 1, reset: now.Add(rl.window)}
		return true, 0
	}
	if w.count >= rl.limit {
		return false, time.Until(w.reset)
	}
	w.count++
	return true, 0
}

// Middleware enforces the rate limit per client IP, returning 429 when exceeded.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok, retry := rl.allow(clientIP(r)); !ok {
			w.Header().Set("Retry-After", strconv.Itoa(int(retry.Seconds())+1))
			response.JSON(w, http.StatusTooManyRequests, map[string]string{
				"error": "too many requests, please slow down",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// clientIP extracts the best-effort client IP, honoring X-Forwarded-For.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
