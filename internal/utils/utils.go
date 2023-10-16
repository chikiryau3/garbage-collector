package utils

import (
	"context"
	"database/sql"
	"github.com/chikiryau3/garbage-collector/internal/configs"
	"github.com/chikiryau3/garbage-collector/internal/logger"
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	pgstorage "github.com/chikiryau3/garbage-collector/internal/pgStorage"
	_ "github.com/jackc/pgx/v5"
)

func ListenForErrors(errs <-chan error, prefix string, logCb func(args ...interface{})) {
	err := <-errs
	if err != nil {
		logCb(prefix, err)
	}
}

func InitPgConnection(conString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", conString)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitMemStorage(config *configs.ServiceConfig, log logger.Logger) metricscollector.Storage {
	storage := memstorage.New(config.StorageConfig)

	if config.Restore {
		err := storage.RestoreFromDump()
		if err != nil {
			log.Error("restore from dump error", err)
		}
	}

	if config.FileStoragePath != "" {
		errs := storage.RunStorageDumper()
		go ListenForErrors(errs, "storage dumper error", log.Error)
	}

	return storage
}

func InitPgStorage(ctx context.Context, db *sql.DB) (pgstorage.PgStorage, error) {
	s := pgstorage.New(db, &pgstorage.Config{})
	err := s.Init(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}
