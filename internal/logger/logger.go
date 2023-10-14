package logger

import (
	"fmt"
	"go.uber.org/zap"
)

type Logger interface {
	Infoln(args ...interface{})
	Error(args ...interface{})
}

func InitLogger() (Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("init logger error %w", err)
	}

	log := *logger.Sugar()

	return &log, nil
}
