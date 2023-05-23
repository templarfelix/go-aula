package http

import (
	"context"
	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"microservice/cmd/infra/config"
	customEchoContext "microservice/cmd/infra/context"
	httpPackage "net/http"
)

func ProvideHttpServer(logger *zap.SugaredLogger) *echo.Echo {
	logger.Info("Executing ProvideHttpServer.")
	echoInstance := echo.New()
	echoInstance.HideBanner = true
	return echoInstance
}

func RegisterHooks(
	lifecycle fx.Lifecycle,
	logger *zap.SugaredLogger,
	config config.Config,
	echoInstance *echo.Echo,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				logger.Info("Listening on ", config.Server.Address)
				go func() {
					if err := echoInstance.Start(config.Server.Address); err != nil && err != httpPackage.ErrServerClosed {
						logger.Fatal("shutting down the server")
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return echoInstance.Shutdown(ctx)
			},
		},
	)
}

func RegisterMiddlewareHooks(
	lifecycle fx.Lifecycle,
	logger *zap.SugaredLogger,
	echoInstance *echo.Echo,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				echoInstance.Use(middleware.RequestID())
				p := prometheus.NewPrometheus("echo", nil)
				p.Use(echoInstance)

				pprof.Register(echoInstance)

				echoInstance.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
					LogURI:    true,
					LogStatus: true,
					LogHost:   true,
					LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
						logger.Info("request",
							zap.String("URI", v.URI),
							zap.Int("status", v.Status),
						)
						return nil
					},
				}))
				echoInstance.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
					return func(c echo.Context) error {
						return next(&customEchoContext.EchoContext{c})
					}
				})
				return nil
			},
			OnStop: func(context.Context) error {
				return nil
			},
		},
	)
}
