package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CalculatorMetrics struct {
	Attempts prometheus.Counter
	Success  prometheus.Counter
	Failed   prometheus.Counter
	Duration prometheus.Histogram
}

type Metrics struct {
	System struct {
		HTTPRequests prometheus.Counter
	}

	ErrorTypes *prometheus.CounterVec

	UI         CalculatorMetrics
	PublicAPI  CalculatorMetrics
	PrivateAPI CalculatorMetrics
}

var M = New()

func New() *Metrics {
	m := &Metrics{}

	// System metrics
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

	// UI metrics
	m.UI.Attempts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tax_web_ui_calc_attempts_total",
		Help: "UI calculation attempts",
	})

	m.UI.Success = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tax_web_ui_calc_success_total",
		Help: "UI calculation success",
	})

	m.UI.Failed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tax_web_ui_calc_failed_total",
		Help: "UI calculation failed",
	})

	m.UI.Duration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "tax_web_ui_calc_duration_seconds",
		Help:    "Duration of UI calculation processing",
		Buckets: prometheus.DefBuckets,
	})

	// Public API metrics
	m.PublicAPI.Attempts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tax_web_public_calc_attempts_total",
		Help: "Public API calculation attempts",
	})

	m.PublicAPI.Success = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tax_web_public_calc_success_total",
		Help: "Public API calculation success",
	})

	m.PublicAPI.Failed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tax_web_public_calc_failed_total",
		Help: "Public API calculation failed",
	})

	m.PublicAPI.Duration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "tax_web_public_calc_duration_seconds",
		Help:    "Duration of public API calculation",
		Buckets: prometheus.DefBuckets,
	})

	// Private API metrics
	m.PrivateAPI.Attempts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tax_web_private_calc_attempts_total",
		Help: "Private API calculation attempts",
	})

	m.PrivateAPI.Success = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tax_web_private_calc_success_total",
		Help: "Private API calculation success",
	})

	m.PrivateAPI.Failed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tax_web_private_calc_failed_total",
		Help: "Private API calculation failed",
	})

	m.PrivateAPI.Duration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "tax_web_private_calc_duration_seconds",
		Help:    "Duration of private API calculation",
		Buckets: prometheus.DefBuckets,
	})

	return m
}
