package metrics_collector

import "github.com/chikiryau3/garbage-collector/internal/mem-storage"

type MetricsCollector interface {
	Gauge(name string, value float64) error
	Count(name string, value int64) error
}

type metricsCollector struct {
	storage mem_storage.MemStorage
}

func New(storage mem_storage.MemStorage) MetricsCollector {
	return &metricsCollector{
		storage: storage,
	}
}
