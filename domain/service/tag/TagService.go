package tag

import (
	"context"
	"microservice/domain/entitie"
	_interface "microservice/domain/interface"
	"time"
)

type tagService struct {
	repo           _interface.TagRepository
	contextTimeout time.Duration
}

func NewTagService(repo _interface.TagRepository, timeout time.Duration) _interface.TagService {
	return &tagService{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (service tagService) GetByID(ctx context.Context, id uint) (entitie.Tag, error) {
	return service.repo.GetByID(ctx, id)
}

func (service tagService) GetAll(ctx context.Context) ([]entitie.Tag, error) {
	return service.repo.GetAll(ctx)
}

func (service tagService) Update(ctx context.Context, ar *entitie.Tag) error {
	return service.repo.Update(ctx, ar)
}

func (service tagService) GetByName(ctx context.Context, title string) (entitie.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()
	res, err := service.repo.GetByName(ctx, title)
	if err != nil {
		return entitie.Tag{}, err
	}
	return res, err
}

func (service tagService) Store(ctx context.Context, data *entitie.Tag) error {
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()
	existedArticle, _ := service.GetByName(ctx, data.Name)
	if existedArticle != (entitie.Tag{}) {
		return _interface.ErrConflict
	}
	return service.repo.Store(ctx, data)

}

func (service tagService) Delete(ctx context.Context, id uint) error {
	return service.repo.Delete(ctx, id)
}
