package repository

import (
	"context"
	"gorm.io/gorm"
	"microservice/domain/entitie"
	_interface "microservice/domain/interface"
)

type tagRepository struct {
	*gorm.DB
}

func NewTagRepository(conn *gorm.DB) _interface.TagRepository {
	return &tagRepository{conn}
}

func (m *tagRepository) Fetch(ctx context.Context, cursor string, num int64) (res []entitie.Tag, nextCursor string, err error) {

	return nil, "", nil
}

func (m *tagRepository) GetByID(ctx context.Context, id int64) (entitie.Tag, error) {
	return entitie.Tag{}, nil
}

func (m *tagRepository) GetByName(ctx context.Context, title string) (entitie.Tag, error) {
	return entitie.Tag{}, nil
}

func (m *tagRepository) Update(ctx context.Context, ar *entitie.Tag) error {
	return nil
}

func (m *tagRepository) Store(ctx context.Context, a *entitie.Tag) error {
	return nil
}

func (m *tagRepository) Delete(ctx context.Context, id int64) error {
	return nil
}
