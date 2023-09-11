package metricscollector

import "github.com/chikiryau3/garbage-collector/internal/memStorage"

type MetricsCollector interface {
	Gauge(name string, value float64) error
	Count(name string, value int64) error
}

type metricsCollector struct {
	storage memstorage.MemStorage
}

func New(storage memstorage.MemStorage) MetricsCollector {
	return &metricsCollector{
		storage: storage,
	}
}
