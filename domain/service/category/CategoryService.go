package category

import (
	"context"
	"time"

	"go.uber.org/zap"
	"microservice/cmd/infra/config"
	"microservice/domain/entitie"
	_interface "microservice/domain/interface"
)

type categoryService struct {
	repo           _interface.CategoryRepository
	logger         *zap.SugaredLogger
	contextTimeout time.Duration
}

func ProvideCategoryService(repo _interface.CategoryRepository, config config.Config, logger *zap.SugaredLogger) _interface.CategoryService {
	return &categoryService{
		repo:           repo,
		logger:         logger,
		contextTimeout: time.Duration(config.Server.Timeout) * time.Second,
	}
}

func (s *categoryService) GetByID(ctx context.Context, id uint) (entitie.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Debugw("failed to get category by id", "id", id, "error", err)
		return entitie.Category{}, err
	}
	return result, nil
}

func (s *categoryService) GetAll(ctx context.Context) ([]entitie.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	result, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Errorw("failed to get all categories", "error", err)
		return nil, err
	}
	return result, nil
}

func (s *categoryService) Update(ctx context.Context, ar *entitie.Category) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if err := s.repo.Update(ctx, ar); err != nil {
		s.logger.Errorw("failed to update category", "id", ar.ID, "error", err)
		return err
	}
	s.logger.Infow("category updated", "id", ar.ID, "name", ar.Name)
	return nil
}

func (s *categoryService) GetByName(ctx context.Context, name string) (entitie.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	result, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return entitie.Category{}, err
	}
	return result, nil
}

func (s *categoryService) Store(ctx context.Context, data *entitie.Category) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	existedEntity, _ := s.GetByName(ctx, data.Name)
	if existedEntity != (entitie.Category{}) {
		s.logger.Warnw("category already exists", "name", data.Name)
		return _interface.ErrConflict
	}

	if err := s.repo.Store(ctx, data); err != nil {
		s.logger.Errorw("failed to store category", "name", data.Name, "error", err)
		return err
	}
	s.logger.Infow("category created", "id", data.ID, "name", data.Name)
	return nil
}

func (s *categoryService) Delete(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorw("failed to delete category", "id", id, "error", err)
		return err
	}
	s.logger.Infow("category deleted", "id", id)
	return nil
}
