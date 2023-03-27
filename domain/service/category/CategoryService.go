package category

import (
	"context"
	"microservice/domain/entitie"
	_interface "microservice/domain/interface"
	"time"
)

type categoryService struct {
	repo           _interface.CategoryRepository
	contextTimeout time.Duration
}

func NewCategoryService(repo _interface.CategoryRepository, timeout time.Duration) _interface.CategoryService {
	return &categoryService{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (service categoryService) GetByID(ctx context.Context, id uint) (entitie.Category, error) {
	return service.repo.GetByID(ctx, id)
}

func (service categoryService) GetAll(ctx context.Context) ([]entitie.Category, error) {
	return service.repo.GetAll(ctx)
}

func (service categoryService) Update(ctx context.Context, ar *entitie.Category) error {
	return service.repo.Update(ctx, ar)
}

func (service categoryService) GetByName(ctx context.Context, title string) (entitie.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()
	res, err := service.repo.GetByName(ctx, title)
	if err != nil {
		return entitie.Category{}, err
	}
	return res, err
}

func (service categoryService) Store(ctx context.Context, data *entitie.Category) error {
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()
	existedArticle, _ := service.GetByName(ctx, data.Name)
	if existedArticle != (entitie.Category{}) {
		return _interface.ErrConflict
	}
	return service.repo.Store(ctx, data)

}

func (service categoryService) Delete(ctx context.Context, id uint) error {
	return service.repo.Delete(ctx, id)
}
