package _interface

import (
	"context"
	"github.com/labstack/echo/v4"
	"microservice/domain/entitie"
)

type TagHandler interface {
	GetByID(echo.Context) error
	Update(echoContext echo.Context) error
	//GetByName(echo.Context) error
	GetAll(echo.Context) error
	Store(echo.Context) error
	Delete(echo.Context) error
}

//go:generate mockery --name TagService
type TagService interface {
	GetByID(ctx context.Context, id uint) (entitie.Tag, error)
	Update(ctx context.Context, ar *entitie.Tag) error
	GetByName(ctx context.Context, title string) (entitie.Tag, error)
	GetAll(ctx context.Context) ([]entitie.Tag, error)
	Store(ctx context.Context, entitie *entitie.Tag) error
	Delete(ctx context.Context, id uint) error
}

//go:generate mockery --name TagRepository
type TagRepository interface {
	GetByID(ctx context.Context, id uint) (entitie.Tag, error)
	GetByName(ctx context.Context, title string) (entitie.Tag, error)
	GetAll(ctx context.Context) ([]entitie.Tag, error)
	Update(ctx context.Context, ar *entitie.Tag) error
	Store(ctx context.Context, a *entitie.Tag) error
	Delete(ctx context.Context, id uint) error
}
