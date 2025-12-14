package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CalculatorMetrics struct {
	Attempts *prometheus.CounterVec
	Success  *prometheus.CounterVec
	Failed   *prometheus.CounterVec
	Duration *prometheus.HistogramVec
}

type Metrics struct {
	System struct {
		HTTPRequests prometheus.Counter
	}

	ErrorTypes *prometheus.CounterVec
	Calculator *CalculatorMetrics
}

var M = New()

func New() *Metrics {
	m := &Metrics{}

	m.Calculator = &CalculatorMetrics{}

	// System
	m.System.HTTPRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tax_web_http_requests_total",
		Help: "Total number of HTTP requests",
	})

	// Errors
	m.ErrorTypes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tax_web_errors_total",
			Help: "Total errors by type",
		},
		[]string{"endpoint", "type"},
	)

	labels := []string{"client", "region"}

	// Attempts
	m.Calculator.Attempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tax_web_calc_attempts_total",
			Help: "Calculation attempts",
		},
		labels,
	)

	// Failed
	m.Calculator.Failed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tax_web_calc_failed_total",
			Help: "Calculation failed",
		},
		labels,
	)

	// Success
	m.Calculator.Success = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tax_web_calc_success_total",
			Help: "Calculation success",
		},
		labels,
	)

	// Duration
	m.Calculator.Duration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tax_web_calc_duration_seconds",
			Help:    "Calculation duration",
			Buckets: prometheus.DefBuckets,
		},
		labels,
	)

	return m
}
