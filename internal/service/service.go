package service

import (
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"go.uber.org/zap"
	"net/http"
)

type Service interface {
	// handlers

	UpdateHandler(w http.ResponseWriter, r *http.Request)
	ValueHandler(w http.ResponseWriter, r *http.Request)
	GetMetricsHTML(w http.ResponseWriter, r *http.Request)

	// middlewares

	WithLogging(next http.Handler) http.Handler
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
