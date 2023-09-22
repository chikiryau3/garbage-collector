package agent

import (
	garbagecollector "github.com/chikiryau3/garbage-collector/internal/clients/garbage-collector"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"time"
)

type Agent interface {
	RunPollChron() error
	RunReporter() error
}

type Config struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
}

type agent struct {
	collector               metricscollector.MetricsCollector
	collectionServiceClient garbagecollector.Client
	config                  Config
}

func New(c metricscollector.MetricsCollector, sc garbagecollector.Client, config Config) Agent {
	return &agent{
		collector:               c,
		collectionServiceClient: sc,
		config:                  config,
	}
}
