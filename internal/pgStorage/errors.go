package pgstorage

import (
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	"github.com/jackc/pgerrcode"
)

func NewPgError(err error) error {
	if pgerrcode.IsConnectionException(err.Error()) {
		return metricscollector.NewStorageRetryableError(err)
	} else {
		return metricscollector.NewStorageError(err)
	}
}
