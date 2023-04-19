package infra

import (
	"go.uber.org/fx"
	"microservice/cmd/infra/config"
	"microservice/cmd/infra/database"
	"microservice/cmd/infra/http"
	"microservice/cmd/infra/log"
	"microservice/cmd/infra/opentelemetry"
)

var Module = fx.Module("infra",

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
	fx.Provide(
		opentelemetry.ProvideTracerProvider,
		opentelemetry.ProvideMeterProvider,
	),

	fx.Invoke(http.RegisterHooks),
	fx.Invoke(http.RegisterMiddlewareHooks),
	fx.Invoke(database.RegisterHooks),
	fx.Invoke(log.RegisterHooks),
	fx.Invoke(opentelemetry.RegisterHooks),
)
