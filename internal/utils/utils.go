package utils

import (
	"context"
	"database/sql"
	pgstorage "github.com/chikiryau3/garbage-collector/internal/pgStorage"
	_ "github.com/jackc/pgx/v5"
)

func ListenForErrors(errs <-chan error, prefix string, logCb func(args ...interface{})) {
	err := <-errs
	if err != nil {
		logCb(prefix, err)
	}
}

func InitPgStorage(ctx context.Context, db *sql.DB) (pgstorage.PgStorage, error) {
	s := pgstorage.New(db, &pgstorage.Config{})
	err := s.Init(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}
