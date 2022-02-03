package monitoring

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "infro"
	subsystem = "secrets-manager"
)

func init() {
	prometheus.MustRegister(
		TotalReqs,
		TotalErrReqs,
		ReqDurationSeconds,
	)
}

// TotalReqs - requests total.
var TotalReqs = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: namespace,
	Subsystem: subsystem,
	Name:      "requests",
	Help:      "Total amount of requests.",
}, []string{"method"})

// TotalErrReqs - request errors total.
var TotalErrReqs = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: namespace,
	Subsystem: subsystem,
	Name:      "request_errors_total",
	Help:      "Total amount of error requests.",
}, []string{"method", "error_code", "status"})

// ReqDurationSeconds - client request duration(seconds).
var ReqDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: namespace,
	Subsystem: subsystem,
	Name:      "request_duration_seconds",
	Help:      "Request duration seconds.",
}, []string{"method", "is_error"})
