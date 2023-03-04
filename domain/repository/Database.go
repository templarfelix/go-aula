package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"microservice/domain/entitie"
)

func Connect(Host, Port, User, Name, Password string) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		Host,
		Port,
		User,
		Name,
		Password)
	return gorm.Open(postgres.Open(connectionString), &gorm.Config{})
}

func Migrate(Instance *gorm.DB) {
	Instance.AutoMigrate(&entitie.Tag{})
	log.Println("Database Migration Completed...")
}
