package repository

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func Connect(Host, Port, User, Name, Password string) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=verify-full",
		Host,
		Port,
		User,
		Name,
		Password)

	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()
	return gorm.Open(postgres.Open(connectionString), &gorm.Config{Logger: logger})
}
