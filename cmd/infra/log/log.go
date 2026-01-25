package log

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	slogger := logger.Sugar()
	slogger.Info("Executing ProvideLogger.")
	return slogger, nil
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
				return logger.Sync()
			},
		},
	)
}
