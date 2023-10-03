package service

import (
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"go.uber.org/zap"
	"net/http"
)

type Service interface {
	// handlers json

	UpdateHandler(w http.ResponseWriter, r *http.Request)
	ValueHandler(w http.ResponseWriter, r *http.Request)

	// handlers plain

	GaugeHandler(w http.ResponseWriter, r *http.Request)
	CounterHandler(w http.ResponseWriter, r *http.Request)
	GetMetric(w http.ResponseWriter, r *http.Request)
	GetMetricsHTML(w http.ResponseWriter, r *http.Request)

	// middlewares

	WithLogging(next http.Handler) http.Handler
}

type MetricData struct {
	mtype string
	name  string
	value any
}

type service struct {
	collector metricscollector.MetricsCollector
	log       zap.SugaredLogger
}

func New(collector metricscollector.MetricsCollector, log zap.SugaredLogger) Service {
	return &service{
		collector: collector,
		log:       log,
	}
}
