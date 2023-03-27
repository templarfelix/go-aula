package main

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	echoContext "microservice/cmd/infra/context"
	"microservice/cmd/infra/env"
	categoryHandler "microservice/domain/handler/category"
	tagHandler "microservice/domain/handler/tag"
	"microservice/domain/repository"
	categoryRepository "microservice/domain/repository/category"
	tagRepository "microservice/domain/repository/tag"
	categoryService "microservice/domain/service/category"
	tagService "microservice/domain/service/tag"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// servidor web echo https://echo.labstack.com
	echoInstance := echo.New()

	// init logger https://github.com/uber-go/zap
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// database https://gorm.io/docs/index.html
	database, err := repository.Connect(env.Config.Database.Host, env.Config.Database.Port, env.Config.Database.User, env.Config.Database.Name, env.Config.Database.Password)
	if err != nil {
		zap.L().Fatal(err.Error(), zap.Error(err))
	}

	// repository
	tagRepo := tagRepository.NewTagRepository(database)
	categoryRepo := categoryRepository.NewCategoryRepository(database)

	// service
	tagService := tagService.NewTagService(tagRepo, 10)
	categoryService := categoryService.NewCategoryService(categoryRepo, 10)

	echoInstance.Use(middleware.RequestID())
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(echoInstance)
	echoInstance.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
	echoInstance.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&echoContext.EchoContext{c})
		}
	})

	// ROUTES
	echoInstance.GET("/debug", func(c echo.Context) error {
		cc := c.(*echoContext.EchoContext)
		cc.Foo()
		cc.Bar()
		return cc.String(200, "OK")
	})

	// handler
	tagHandler.NewTagHandler(echoInstance, tagService)
	categoryHandler.NewCategoryHandler(echoInstance, categoryService)

	// Start server
	go func() {
		if err := echoInstance.Start(env.Config.Server.Address); err != nil && err != http.ErrServerClosed {
			echoInstance.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(env.Config.Server.Timeout)*time.Second)
	defer cancel()
	if err := echoInstance.Shutdown(ctx); err != nil {
		echoInstance.Logger.Fatal(err)
	}
}
