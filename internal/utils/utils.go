package utils

import (
	"database/sql"
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
