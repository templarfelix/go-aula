package tag

import (
	"context"
	"time"

	"go.uber.org/zap"
	"microservice/cmd/infra/config"
	"microservice/domain/entitie"
	_interface "microservice/domain/interface"
)

type tagService struct {
	repo           _interface.TagRepository
	logger         *zap.SugaredLogger
	contextTimeout time.Duration
}

func ProvideTagService(repo _interface.TagRepository, config config.Config, logger *zap.SugaredLogger) _interface.TagService {
	return &tagService{
		repo:           repo,
		logger:         logger,
		contextTimeout: time.Duration(config.Server.Timeout) * time.Second,
	}
}

func (s *tagService) GetByID(ctx context.Context, id uint) (entitie.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Debugw("failed to get tag by id", "id", id, "error", err)
		return entitie.Tag{}, err
	}
	return result, nil
}

func (s *tagService) GetAll(ctx context.Context) ([]entitie.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	result, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Errorw("failed to get all tags", "error", err)
		return nil, err
	}
	return result, nil
}

func (s *tagService) Update(ctx context.Context, ar *entitie.Tag) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if err := s.repo.Update(ctx, ar); err != nil {
		s.logger.Errorw("failed to update tag", "id", ar.ID, "error", err)
		return err
	}
	s.logger.Infow("tag updated", "id", ar.ID, "name", ar.Name)
	return nil
}

func (s *tagService) GetByName(ctx context.Context, name string) (entitie.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	result, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return entitie.Tag{}, err
	}
	return result, nil
}

func (s *tagService) Store(ctx context.Context, data *entitie.Tag) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	existedEntity, _ := s.GetByName(ctx, data.Name)
	if existedEntity != (entitie.Tag{}) {
		s.logger.Warnw("tag already exists", "name", data.Name)
		return _interface.ErrConflict
	}

	if err := s.repo.Store(ctx, data); err != nil {
		s.logger.Errorw("failed to store tag", "name", data.Name, "error", err)
		return err
	}
	s.logger.Infow("tag created", "id", data.ID, "name", data.Name)
	return nil
}

func (s *tagService) Delete(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorw("failed to delete tag", "id", id, "error", err)
		return err
	}
	s.logger.Infow("tag deleted", "id", id)
	return nil
}
