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

type agent struct {
	collector               metricscollector.MetricsCollector
	collectionServiceClient garbagecollector.Client
	pollInterval            time.Duration
	reportInterval          time.Duration
}

func New(
	collector metricscollector.MetricsCollector,
	collectionServiceClient garbagecollector.Client,
	pollInterval time.Duration,
	reportInterval time.Duration,
) Agent {
	return &agent{
		collector:               collector,
		collectionServiceClient: collectionServiceClient,
		pollInterval:            pollInterval,
		reportInterval:          reportInterval,
	}
}
