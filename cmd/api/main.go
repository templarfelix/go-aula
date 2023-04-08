package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"microservice/cmd/infra/config"
	"microservice/cmd/infra/database"
	"microservice/cmd/infra/http"
	"microservice/cmd/infra/log"
	tagHandler "microservice/domain/handler/tag"
	tagRepository "microservice/domain/repository/tag"
	tagService "microservice/domain/service/tag"
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
		fx.Provide(tagRepository.ProvideTagRepository),
		fx.Provide(tagService.ProvideTagService),
		fx.Provide(tagHandler.ProvideTagHandler),
		fx.Invoke(http.RegisterHooks),
		fx.Invoke(http.RegisterMiddlewareHooks),
		fx.Invoke(database.RegisterHooks),
		fx.Invoke(tagHandler.RegisterHooks),
	).Run()

}
