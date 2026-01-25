package tag_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"microservice/cmd/infra/config"
	"microservice/domain/entitie"
	_interface "microservice/domain/interface"
	"microservice/domain/interface/mocks"
	"microservice/domain/service/tag"
)

func setupService(mockRepo *mocks.TagRepository) _interface.TagService {
	logger := zap.NewNop().Sugar()
	cfg := config.Config{
		Server: config.ServerConfig{Timeout: 30},
	}
	return tag.ProvideTagService(mockRepo, cfg, logger)
}

func TestGetByID(t *testing.T) {
	mockRepo := new(mocks.TagRepository)
	mockTag := entitie.Tag{
		Name: "Hello",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(mockTag, nil).Once()
		service := setupService(mockRepo)

		result, err := service.GetByID(context.TODO(), mockTag.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockTag.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(entitie.Tag{}, errors.New("unexpected")).Once()
		service := setupService(mockRepo)

		result, err := service.GetByID(context.TODO(), mockTag.ID)

		assert.Error(t, err)
		assert.Equal(t, entitie.Tag{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAll(t *testing.T) {
	mockRepo := new(mocks.TagRepository)

	t.Run("success", func(t *testing.T) {
		mockTags := []entitie.Tag{
			{Name: "Tag 1"},
			{Name: "Tag 2"},
		}
		mockRepo.On("GetAll", mock.Anything).Return(mockTags, nil).Once()
		service := setupService(mockRepo)

		result, err := service.GetAll(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetAll", mock.Anything).Return([]entitie.Tag{}, errors.New("database error")).Once()
		service := setupService(mockRepo)

		result, err := service.GetAll(context.TODO())

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	mockRepo := new(mocks.TagRepository)
	mockTag := entitie.Tag{
		Name: "Hello",
	}

	t.Run("success", func(t *testing.T) {
		tempTag := mockTag
		tempTag.ID = 0

		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(entitie.Tag{}, _interface.ErrNotFound).Once()
		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*entitie.Tag")).Return(nil).Once()
		service := setupService(mockRepo)

		err := service.Store(context.TODO(), &tempTag)

		assert.NoError(t, err)
		assert.Equal(t, mockTag.Name, tempTag.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("existing_name", func(t *testing.T) {
		existingTag := mockTag

		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(existingTag, nil).Once()
		service := setupService(mockRepo)

		err := service.Store(context.TODO(), &mockTag)

		assert.Error(t, err)
		assert.Equal(t, _interface.ErrConflict, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("store_error", func(t *testing.T) {
		tempTag := mockTag
		tempTag.ID = 0

		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(entitie.Tag{}, _interface.ErrNotFound).Once()
		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*entitie.Tag")).Return(errors.New("database error")).Once()
		service := setupService(mockRepo)

		err := service.Store(context.TODO(), &tempTag)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockRepo := new(mocks.TagRepository)

	t.Run("success", func(t *testing.T) {
		mockTag := entitie.Tag{Name: "Updated"}
		mockTag.ID = 1
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entitie.Tag")).Return(nil).Once()
		service := setupService(mockRepo)

		err := service.Update(context.TODO(), &mockTag)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockTag := entitie.Tag{Name: "Updated"}
		mockTag.ID = 1
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entitie.Tag")).Return(errors.New("database error")).Once()
		service := setupService(mockRepo)

		err := service.Update(context.TODO(), &mockTag)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockRepo := new(mocks.TagRepository)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(nil).Once()
		service := setupService(mockRepo)

		err := service.Delete(context.TODO(), uint(1))

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(errors.New("database error")).Once()
		service := setupService(mockRepo)

		err := service.Delete(context.TODO(), uint(1))

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
