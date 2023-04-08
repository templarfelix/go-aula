package log

import (
	"go.uber.org/zap"
)

func ProvideLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	slogger := logger.Sugar()
	slogger.Info("Executing ProvideLogger.")
	return slogger
}
