package web

import (
	"net/http"
	"sync"
	"time"
)

// SecurityHeaders sets defensive HTTP response headers on every reply.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// Leaflet needs an inline <script> bootstrap on the verdict page,
		// and OSM tile servers are an external img + connect source.
		h.Set("Content-Security-Policy",
			"default-src 'self'; "+
				"img-src 'self' data: https://*.tile.openstreetmap.org https:; "+
				"script-src 'self' 'unsafe-inline'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"connect-src 'self' https://*.tile.openstreetmap.org; "+
				"frame-ancestors 'none'; "+
				"base-uri 'self'")
		next.ServeHTTP(w, r)
	})
}

// IPRateLimiter is a simple per-IP token-bucket limiter for hot
// endpoints (e.g. /verify).
type IPRateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    int           // tokens per window
	window  time.Duration // refill window
}

type bucket struct {
	tokens int
	resets time.Time
}

// NewIPRateLimiter allows `rate` requests per `window` per IP.
func NewIPRateLimiter(rate int, window time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		buckets: make(map[string]*bucket),
		rate:    rate,
		window:  window,
	}
}

// Middleware applies the limiter; replies 429 when exhausted.
func (l *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		l.mu.Lock()
		b, ok := l.buckets[ip]
		now := time.Now()
		if !ok || now.After(b.resets) {
			b = &bucket{tokens: l.rate, resets: now.Add(l.window)}
			l.buckets[ip] = b
		}
		if b.tokens <= 0 {
			retry := int(time.Until(b.resets).Seconds()) + 1
			l.mu.Unlock()
			w.Header().Set("Retry-After", strInt(retry))
			http.Error(w, "rate limit", http.StatusTooManyRequests)
			return
		}
		b.tokens--
		l.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		return v
	}
	return r.RemoteAddr
}

func strInt(n int) string {
	if n <= 0 {
		return "1"
	}
	const digits = "0123456789"
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{digits[n%10]}, buf...)
		n /= 10
	}
	return string(buf)
}
