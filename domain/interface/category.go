package _interface

import (
	"context"
	"github.com/labstack/echo/v4"
	"microservice/domain/entitie"
)

type CategoryHandler interface {
	GetByID(echo.Context) error
	Update(echoContext echo.Context) error
	//GetByName(echo.Context) error
	GetAll(echo.Context) error
	Store(echo.Context) error
	Delete(echo.Context) error
}

//go:generate mockery --name CategoryService
type CategoryService interface {
	GetByID(ctx context.Context, id uint) (entitie.Category, error)
	Update(ctx context.Context, ar *entitie.Category) error
	GetByName(ctx context.Context, title string) (entitie.Category, error)
	GetAll(ctx context.Context) ([]entitie.Category, error)
	Store(context.Context, *entitie.Category) error
	Delete(ctx context.Context, id uint) error
}

//go:generate mockery --name CategoryRepository
type CategoryRepository interface {
	GetByID(ctx context.Context, id uint) (entitie.Category, error)
	GetByName(ctx context.Context, title string) (entitie.Category, error)
	GetAll(ctx context.Context) ([]entitie.Category, error)
	Update(ctx context.Context, ar *entitie.Category) error
	Store(ctx context.Context, a *entitie.Category) error
	Delete(ctx context.Context, id uint) error
}
