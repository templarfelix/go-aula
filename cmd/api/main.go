package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"microservice/cmd/infra/config"
	"microservice/cmd/infra/database"
	"microservice/cmd/infra/http"
	"microservice/cmd/infra/log"
	"microservice/domain/handler"
	"microservice/domain/repository"
	"microservice/domain/service"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.SugaredLogger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Desugar()}
		}),
		fx.Provide(
			http.ProvideHttpServer,
		),
		fx.Provide(
			config.ProvideConfig,
		),
		fx.Provide(
			log.ProvideLogger,
		),
		fx.Provide(
			database.ProvideDatabase,
		),

		fx.Invoke(http.RegisterHooks),
		fx.Invoke(http.RegisterMiddlewareHooks),
		fx.Invoke(database.RegisterHooks),

		repository.Module,
		service.Module,
		handler.Module,
	).Run()

}
