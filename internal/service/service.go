package service

import (
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"go.uber.org/zap"
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
	WithLogging(next http.Handler) http.Handler
}

type MetricData struct {
	mtype string
	name  string
	value any
}

type mKey struct{}

type service struct {
	collector            metricscollector.MetricsCollector
	log                  zap.SugaredLogger
	metricDataContextKey struct{}
}

func New(collector metricscollector.MetricsCollector, log zap.SugaredLogger) Service {
	return &service{
		collector:            collector,
		log:                  log,
		metricDataContextKey: mKey{},
	}
}
