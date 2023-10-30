package agent

import (
	garbagecollector "github.com/chikiryau3/garbage-collector/internal/clients/garbage-collector"
	"github.com/chikiryau3/garbage-collector/internal/logger"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"time"
)

type Config struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
}

type agent struct {
	collector               metricscollector.MetricsCollector
	collectionServiceClient garbagecollector.Client
	config                  *Config
	log                     logger.Logger
}

func New(c metricscollector.MetricsCollector, sc garbagecollector.Client, config *Config, log logger.Logger) *agent {
	return &agent{
		collector:               c,
		collectionServiceClient: sc,
		config:                  config,
		log:                     log,
	}
}
