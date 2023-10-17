package common

import (
	"context"
	"database/sql"
	"github.com/chikiryau3/garbage-collector/internal/configs"
	"github.com/chikiryau3/garbage-collector/internal/logger"
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
	metricscollector "github.com/chikiryau3/garbage-collector/internal/metricsCollector"
	pgstorage "github.com/chikiryau3/garbage-collector/internal/pgStorage"
	"github.com/chikiryau3/garbage-collector/internal/utils"
)

func InitMemStorage(config *configs.ServiceConfig, log logger.Logger) metricscollector.Storage {
	s := memstorage.New(config.StorageConfig)

	if config.Restore {
		err := s.RestoreFromDump()
		if err != nil {
			log.Error("restore from dump error", err)
		}
	}

	if config.FileStoragePath != "" {
		errs := s.RunStorageDumper()
		go utils.ListenForErrors(errs, "storage dumper error", log.Error)
	}

	return s
}

func InitPgConnection(conString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", conString)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)

	return db, nil
}

func InitPgStorage(ctx context.Context, db *sql.DB) (pgstorage.PgStorage, error) {
	s := pgstorage.New(db, &pgstorage.Config{})
	err := s.Init(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}
