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

	WithMetricData(next http.Handler) http.Handler
}

type MetricData struct {
	mtype string
	name  string
	value any
}

type mKey struct{}

type service struct {
	collector            metricscollector.MetricsCollector
	metricDataContextKey struct{}
}

func New(collector metricscollector.MetricsCollector) Service {
	return &service{
		collector:            collector,
		metricDataContextKey: mKey{},
	}
}
