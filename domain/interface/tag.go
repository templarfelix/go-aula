package _interface

import (
	"context"
	"microservice/domain/entitie"
)

// TagUsecase represent the Tag's usecases
//
//go:generate mockery --name TagService
type TagService interface {
	GetByID(ctx context.Context, id int64) (entitie.Tag, error)
	Update(ctx context.Context, ar *entitie.Tag) error
	GetByName(ctx context.Context, title string) (entitie.Tag, error)
	Store(context.Context, *entitie.Tag) error
	Delete(ctx context.Context, id int64) error
}

// TagRepository represent the Tag's repository contract
//
//go:generate mockery --name TagRepository
type TagRepository interface {
	GetByID(ctx context.Context, id int64) (entitie.Tag, error)
	GetByName(ctx context.Context, title string) (entitie.Tag, error)
	Update(ctx context.Context, ar *entitie.Tag) error
	Store(ctx context.Context, a *entitie.Tag) error
	Delete(ctx context.Context, id int64) error
}
