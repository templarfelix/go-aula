package database

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"microservice/cmd/infra/config"
	"microservice/domain/entitie"
	"moul.io/zapgorm2"
)

func ProvideDatabase(logger *zap.SugaredLogger, config config.Config) *gorm.DB {

	logger.Info("Executing ProvideHttpServer.")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=verify-full",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Name,
		config.Database.Password)

	loggerZapgorm2 := zapgorm2.New(logger.Desugar())
	loggerZapgorm2.SetAsDefault()

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{Logger: loggerZapgorm2})
	if err != nil {
		logger.Fatal("error on start database", zap.Error(err))
	}
	return db

}

func RegisterHooks(
	lifecycle fx.Lifecycle,
	logger *zap.SugaredLogger,
	db *gorm.DB,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				connection, err := db.DB()
				if err != nil {
					logger.Fatal("error on start db session")
				}
				db.AutoMigrate(&entitie.Tag{})
				db.AutoMigrate(&entitie.Category{})
				return connection.Ping()
			},
			OnStop: func(context.Context) error {
				connection, err := db.DB()
				if err != nil {
					logger.Fatal("shutting down the database")
				}
				return connection.Close()
			},
		},
	)
}
