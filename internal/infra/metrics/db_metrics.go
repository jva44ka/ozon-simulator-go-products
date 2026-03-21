package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type DbMetrics struct {
	requestsTotal          *prometheus.CounterVec
	optimisticLockFailures prometheus.Counter
}

func NewDbMetrics() *DbMetrics {
	return &DbMetrics{
		requestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "products_db_requests_total",
			Help: "Total number of products DB requests",
		}, []string{"method", "status"}),
		optimisticLockFailures: promauto.NewCounter(prometheus.CounterOpts{
			Name: "products_db_optimistic_lock_failures_total",
			Help: "Total number of optimistic lock failures in product count updates",
		}),
	}
}

func (rm *DbMetrics) ReportRequest(method string, status string) {
	rm.requestsTotal.WithLabelValues(method, status).Inc()
}

func (rm *DbMetrics) ReportOptimisticLockFailure() {
	rm.optimisticLockFailures.Inc()
}
