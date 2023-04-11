package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"microservice/cmd/infra"
	"microservice/domain/handler"
	"microservice/domain/repository"
	"microservice/domain/service"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.SugaredLogger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Desugar()}
		}),
		infra.Module,
		repository.Module,
		service.Module,
		handler.Module,
	).Run()

}
