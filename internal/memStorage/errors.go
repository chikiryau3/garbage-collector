package memstorage

import (
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
)

func IsRetriable(_ error) bool {
	// todo: find out how to detect race condition error
	return false
}

func NewMemStorageError(err error) error {
	if IsRetriable(err) {
		return metricscollector.NewStorageRetryableError(err)
	} else {
		return metricscollector.NewStorageError(err)
	}
}
