package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler exposes Joblantern's Prometheus metrics on /metrics
// and provides middleware that records request count + latency.
type MetricsHandler struct {
	reqCount   *prometheus.CounterVec
	reqLatency *prometheus.HistogramVec
}

// NewMetrics registers metrics and returns the handler.
func NewMetrics() *MetricsHandler {
	reqCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "joblantern_http_requests_total",
			Help: "Total HTTP requests partitioned by method, path, status.",
		},
		[]string{"method", "path", "status"},
	)
	reqLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "joblantern_http_request_duration_seconds",
			Help:    "HTTP request latency partitioned by method, path.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	prometheus.MustRegister(reqCount, reqLatency)
	return &MetricsHandler{reqCount: reqCount, reqLatency: reqLatency}
}

// Endpoint returns the /metrics http.Handler.
func (m *MetricsHandler) Endpoint() http.Handler {
	return promhttp.Handler()
}

// Middleware wraps an http.Handler to record metrics.
func (m *MetricsHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		path := r.URL.Path
		// Don't explode the cardinality on UUIDs.
		if len(path) > 80 {
			path = path[:80]
		}
		m.reqCount.WithLabelValues(r.Method, path, strconv.Itoa(sw.status)).Inc()
		m.reqLatency.WithLabelValues(r.Method, path).Observe(time.Since(start).Seconds())
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
