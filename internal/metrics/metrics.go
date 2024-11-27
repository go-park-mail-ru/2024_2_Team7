package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response times for handler in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "service", "status_code"},
	)

	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of requests",
		},
		[]string{"path", "method", "service", "status_code"},
	)

	ErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of error requests",
		},
		[]string{"path", "method", "service", "status_code"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(ErrorCount)
}
