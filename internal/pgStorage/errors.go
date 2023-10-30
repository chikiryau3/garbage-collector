package pgstorage

import (
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"github.com/jackc/pgerrcode"
)

func IsRetryable(errString string) bool {
	return pgerrcode.IsInsufficientResources(errString) || pgerrcode.IsConnectionException(errString)
}

func NewPgError(err error) error {
	if IsRetryable(err.Error()) {
		return metricscollector.NewStorageRetryableError(err)
	} else {
		return metricscollector.NewStorageError(err)
	}
}
