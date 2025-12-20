package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// HTTP request metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	httpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "endpoint"},
	)

	httpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "endpoint"},
	)
)

// InitPrometheus registers all Prometheus metrics
func InitPrometheus() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestSize)
	prometheus.MustRegister(httpResponseSize)
}

// PrometheusMiddleware returns a Gin middleware for Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Record request size
		requestSize := c.Request.ContentLength
		if requestSize == -1 {
			requestSize = 0
		}

		c.Next()

		// Record metrics after request is processed
		duration := time.Since(startTime).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		httpRequestsTotal.WithLabelValues(c.Request.Method, endpoint, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, endpoint).Observe(duration)
		httpRequestSize.WithLabelValues(c.Request.Method, endpoint).Observe(float64(requestSize))
		httpResponseSize.WithLabelValues(c.Request.Method, endpoint).Observe(float64(c.Writer.Size()))
	}
}
