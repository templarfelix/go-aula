package _interface

import (
	"context"
	"microservice/domain/entitie"
)

//go:generate mockery --name TagService
type TagService interface {
	GetByID(ctx context.Context, id uint) (entitie.Tag, error)
	Update(ctx context.Context, ar *entitie.Tag) error
	GetByName(ctx context.Context, title string) (entitie.Tag, error)
	GetAll(ctx context.Context) ([]entitie.Tag, error)
	Store(context.Context, *entitie.Tag) error
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
