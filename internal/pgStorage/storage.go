package pgStorage

import (
	"context"
	"database/sql"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
)

type PgStorage metricscollector.Storage

type Config struct {
}

type storage struct {
	db     *sql.DB
	config *Config
}

func New(db *sql.DB, c *Config) PgStorage {
	return &storage{
		db:     db,
		config: c,
	}
}

func (s *storage) CheckConnection(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *storage) WriteMetric(name string, value any) error {
	return nil
}

func (s *storage) ReadMetric(name string) (any, bool) {
	return nil, true
}

func (s *storage) GetData() (*metricscollector.StorageData, error) {
	return nil, nil
}

func (s *storage) RunStorageDumper() <-chan error {
	return nil
}

func (s *storage) RestoreFromDump() error {
	return nil
}
