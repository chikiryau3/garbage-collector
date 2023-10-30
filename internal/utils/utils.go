package utils

import (
	_ "github.com/jackc/pgx/v5"
)

func ListenForErrors(errs <-chan error, prefix string, logCb func(args ...interface{})) {
	err := <-errs
	if err != nil {
		logCb(prefix, err)
	}
}
