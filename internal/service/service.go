package service

import (
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"net/http"
)

type Service interface {
	// handlers

	GaugeHandler(w http.ResponseWriter, r *http.Request)
	CounterHandler(w http.ResponseWriter, r *http.Request)
	GetMetric(w http.ResponseWriter, r *http.Request)
	GetMetricsHTML(w http.ResponseWriter, r *http.Request)

	// middlewares

	//MetricsCtx(next http.Handler) http.Handler
	WithMetricName(next http.Handler) http.Handler
	WithMetricValue(next http.Handler) http.Handler

	// formatters

	FormatGaugeInput(metricNameRaw any, metricValueRaw any) (string, float64, error)
	FormatCounterInput(metricNameRaw any, metricValueRaw any) (string, int64, error)
}

type service struct {
	collector metricscollector.MetricsCollector
}

func New(collector metricscollector.MetricsCollector) Service {
	return &service{
		collector: collector,
	}
}
