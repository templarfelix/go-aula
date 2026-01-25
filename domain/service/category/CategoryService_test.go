package category_test

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
	"microservice/domain/service/category"
)

func setupService(mockRepo *mocks.CategoryRepository) _interface.CategoryService {
	logger := zap.NewNop().Sugar()
	cfg := config.Config{
		Server: config.ServerConfig{Timeout: 30},
	}
	return category.ProvideCategoryService(mockRepo, cfg, logger)
}

func TestGetByID(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	mockCategory := entitie.Category{
		Name: "Hello",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(mockCategory, nil).Once()
		service := setupService(mockRepo)

		result, err := service.GetByID(context.TODO(), mockCategory.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockCategory.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(entitie.Category{}, errors.New("unexpected")).Once()
		service := setupService(mockRepo)

		result, err := service.GetByID(context.TODO(), mockCategory.ID)

		assert.Error(t, err)
		assert.Equal(t, entitie.Category{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAll(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)

	t.Run("success", func(t *testing.T) {
		mockCategories := []entitie.Category{
			{Name: "Category 1"},
			{Name: "Category 2"},
		}
		mockRepo.On("GetAll", mock.Anything).Return(mockCategories, nil).Once()
		service := setupService(mockRepo)

		result, err := service.GetAll(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty", func(t *testing.T) {
		mockRepo.On("GetAll", mock.Anything).Return([]entitie.Category{}, nil).Once()
		service := setupService(mockRepo)

		result, err := service.GetAll(context.TODO())

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetAll", mock.Anything).Return([]entitie.Category{}, errors.New("database error")).Once()
		service := setupService(mockRepo)

		result, err := service.GetAll(context.TODO())

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	mockCategory := entitie.Category{
		Name: "Hello",
	}

	t.Run("success", func(t *testing.T) {
		tempCategory := mockCategory
		tempCategory.ID = 0

		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(entitie.Category{}, _interface.ErrNotFound).Once()
		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*entitie.Category")).Return(nil).Once()
		service := setupService(mockRepo)

		err := service.Store(context.TODO(), &tempCategory)

		assert.NoError(t, err)
		assert.Equal(t, mockCategory.Name, tempCategory.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("existing_name", func(t *testing.T) {
		existingCategory := mockCategory

		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(existingCategory, nil).Once()
		service := setupService(mockRepo)

		err := service.Store(context.TODO(), &mockCategory)

		assert.Error(t, err)
		assert.Equal(t, _interface.ErrConflict, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("store_error", func(t *testing.T) {
		tempCategory := mockCategory
		tempCategory.ID = 0

		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(entitie.Category{}, _interface.ErrNotFound).Once()
		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*entitie.Category")).Return(errors.New("database error")).Once()
		service := setupService(mockRepo)

		err := service.Store(context.TODO(), &tempCategory)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)

	t.Run("success", func(t *testing.T) {
		mockCategory := entitie.Category{Name: "Updated"}
		mockCategory.ID = 1
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entitie.Category")).Return(nil).Once()
		service := setupService(mockRepo)

		err := service.Update(context.TODO(), &mockCategory)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockCategory := entitie.Category{Name: "Updated"}
		mockCategory.ID = 1
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entitie.Category")).Return(errors.New("database error")).Once()
		service := setupService(mockRepo)

		err := service.Update(context.TODO(), &mockCategory)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)

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

func TestGetByName(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)

	t.Run("success", func(t *testing.T) {
		mockCategory := entitie.Category{Name: "Test"}
		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(mockCategory, nil).Once()
		service := setupService(mockRepo)

		result, err := service.GetByName(context.TODO(), "Test")

		assert.NoError(t, err)
		assert.Equal(t, mockCategory.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(entitie.Category{}, _interface.ErrNotFound).Once()
		service := setupService(mockRepo)

		result, err := service.GetByName(context.TODO(), "NonExistent")

		assert.Error(t, err)
		assert.Equal(t, entitie.Category{}, result)
		mockRepo.AssertExpectations(t)
	})
}
