package main

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	customContext "microservice/cmd/infra/context"
	"microservice/cmd/infra/env"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	e := echo.New()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	e.Use(middleware.RequestID())
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&customContext.CustomContext{c})
		}
	})

	// ROUTES

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/debug", func(c echo.Context) error {
		cc := c.(*customContext.CustomContext)
		cc.Foo()
		cc.Bar()
		return cc.String(200, "OK")
	})

	// Start server
	go func() {
		if err := e.Start(env.Config.Server.Address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(env.Config.Server.Timeout)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
