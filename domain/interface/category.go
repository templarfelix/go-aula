package _interface

import (
	"context"
	"microservice/domain/entitie"
)

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
