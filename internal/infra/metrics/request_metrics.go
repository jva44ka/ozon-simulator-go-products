package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type RequestMetrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewRequestMetrics() *RequestMetrics {
	return &RequestMetrics{
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "products_grpc_requests_total",
				Help: "Total number of gRPC requests",
			},
			[]string{"method", "code"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "products_grpc_request_duration_seconds",
				Help:    "Duration of gRPC requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method"},
		),
	}
}

func (rm *RequestMetrics) ReportRequestInfo(methodName string, code string, duration time.Duration) {
	rm.requestsTotal.WithLabelValues(methodName, code).Inc()
	rm.requestDuration.WithLabelValues(methodName).Observe(duration.Seconds())
}
