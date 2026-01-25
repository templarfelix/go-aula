package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"microservice/cmd/infra/config"
	customEchoContext "microservice/cmd/infra/context"
)

const (
	defaultShutdownTimeout = 30 * time.Second
	defaultReadTimeout     = 15 * time.Second
	defaultWriteTimeout    = 15 * time.Second
	defaultIdleTimeout     = 60 * time.Second
)

func ProvideHttpServer(logger *zap.SugaredLogger) *echo.Echo {
	logger.Info("Executing ProvideHttpServer.")
	echoInstance := echo.New()
	echoInstance.HideBanner = true
	echoInstance.HidePort = true

	// Configure timeouts for the underlying HTTP server
	echoInstance.Server.ReadTimeout = defaultReadTimeout
	echoInstance.Server.WriteTimeout = defaultWriteTimeout
	echoInstance.Server.IdleTimeout = defaultIdleTimeout

	return echoInstance
}

func RegisterHooks(
	lifecycle fx.Lifecycle,
	shutdowner fx.Shutdowner,
	logger *zap.SugaredLogger,
	config config.Config,
	echoInstance *echo.Echo,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Infow("Starting HTTP server", "address", config.Server.Address)
				go func() {
					if err := echoInstance.Start(config.Server.Address); err != nil && err != http.ErrServerClosed {
						logger.Errorw("HTTP server error", "error", err)
						// Signal fx to shut down gracefully
						if err := shutdowner.Shutdown(); err != nil {
							logger.Errorw("Error during shutdown", "error", err)
						}
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("Shutting down HTTP server...")

				// Create a context with timeout for graceful shutdown
				shutdownCtx, cancel := context.WithTimeout(ctx, defaultShutdownTimeout)
				defer cancel()

				if err := echoInstance.Shutdown(shutdownCtx); err != nil {
					logger.Errorw("Error during HTTP server shutdown", "error", err)
					return fmt.Errorf("http server shutdown error: %w", err)
				}

				logger.Info("HTTP server stopped gracefully")
				return nil
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
				// Recovery middleware - recover from panics
				echoInstance.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
					StackSize: 1 << 10, // 1 KB
					LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
						logger.Errorw("Panic recovered",
							"error", err,
							"stack", string(stack),
							"uri", c.Request().RequestURI,
						)
						return nil
					},
				}))

				// Request ID middleware
				echoInstance.Use(middleware.RequestID())

				// Prometheus metrics
				p := prometheus.NewPrometheus("echo", nil)
				p.Use(echoInstance)

				// pprof for profiling
				pprof.Register(echoInstance)

				// Request logging
				echoInstance.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
					LogURI:      true,
					LogStatus:   true,
					LogHost:     true,
					LogMethod:   true,
					LogLatency:  true,
					LogRemoteIP: true,
					LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
						logger.Infow("request",
							"method", v.Method,
							"uri", v.URI,
							"status", v.Status,
							"latency", v.Latency.String(),
							"remote_ip", v.RemoteIP,
							"request_id", c.Response().Header().Get(echo.HeaderXRequestID),
						)
						return nil
					},
				}))

				// Custom context middleware
				echoInstance.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
					return func(c echo.Context) error {
						return next(&customEchoContext.EchoContext{Context: c})
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
