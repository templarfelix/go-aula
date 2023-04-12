package log

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	slogger := logger.Sugar()
	slogger.Info("Executing ProvideLogger.")
	return slogger
}

func RegisterHooks(
	lifecycle fx.Lifecycle,
	logger *zap.SugaredLogger,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				return nil
			},
			OnStop: func(context.Context) error {
				//return logger.Desugar().Sync()
				return nil
			},
		},
	)
}
