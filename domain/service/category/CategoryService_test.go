package category_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"microservice/domain/entitie"
	_interface "microservice/domain/interface"
	"microservice/domain/interface/mocks"
	"microservice/domain/service/category"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	mockArticleRepo := new(mocks.CategoryRepository)
	mockArticle := entitie.Category{
		Name: "Hello",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(mockArticle, nil).Once()
		u := category.NewCategoryService(mockArticleRepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockArticleRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(entitie.Category{}, errors.New("Unexpected")).Once()

		u := category.NewCategoryService(mockArticleRepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		assert.Equal(t, entitie.Category{}, a)

		mockArticleRepo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	mockArticleRepo := new(mocks.CategoryRepository)
	mockArticle := entitie.Category{
		Name: "Hello",
	}

	t.Run("success", func(t *testing.T) {
		tempMockArticle := mockArticle
		tempMockArticle.ID = 0

		mockArticleRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(entitie.Category{}, _interface.ErrNotFound).Once()
		mockArticleRepo.On("Store", mock.Anything, mock.AnythingOfType("*entitie.Category")).Return(nil).Once()

		u := category.NewCategoryService(mockArticleRepo, time.Second*2)

		err := u.Store(context.TODO(), &tempMockArticle)

		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Name, tempMockArticle.Name)
		mockArticleRepo.AssertExpectations(t)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingArticle := mockArticle

		mockArticleRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(existingArticle, nil).Once()

		u := category.NewCategoryService(mockArticleRepo, time.Second*2)

		err := u.Store(context.TODO(), &mockArticle)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
	})
}
