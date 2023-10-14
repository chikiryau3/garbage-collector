package service

import (
	"github.com/chikiryau3/garbage-collector/internal/logger"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
)

type MetricData struct {
	mtype string
	name  string
	value string
}

type service struct {
	collector metricscollector.MetricsCollector
	log       logger.Logger
}

func New(collector metricscollector.MetricsCollector, log logger.Logger) *service {
	return &service{
		collector: collector,
		log:       log,
	}
}
