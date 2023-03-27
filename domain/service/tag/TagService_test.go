package tag_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"microservice/domain/entitie"
	_interface "microservice/domain/interface"
	"microservice/domain/interface/mocks"
	"microservice/domain/service/tag"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	mockArticleRepo := new(mocks.TagRepository)
	mockArticle := entitie.Tag{
		Name: "Hello",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(mockArticle, nil).Once()
		u := tag.NewTagService(mockArticleRepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockArticleRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(entitie.Tag{}, errors.New("Unexpected")).Once()

		u := tag.NewTagService(mockArticleRepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		assert.Equal(t, entitie.Tag{}, a)

		mockArticleRepo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	mockArticleRepo := new(mocks.TagRepository)
	mockArticle := entitie.Tag{
		Name: "Hello",
	}

	t.Run("success", func(t *testing.T) {
		tempMockArticle := mockArticle
		tempMockArticle.ID = 0

		mockArticleRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(entitie.Tag{}, _interface.ErrNotFound).Once()
		mockArticleRepo.On("Store", mock.Anything, mock.AnythingOfType("*entitie.Tag")).Return(nil).Once()

		u := tag.NewTagService(mockArticleRepo, time.Second*2)

		err := u.Store(context.TODO(), &tempMockArticle)

		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Name, tempMockArticle.Name)
		mockArticleRepo.AssertExpectations(t)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingArticle := mockArticle

		mockArticleRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(existingArticle, nil).Once()

		u := tag.NewTagService(mockArticleRepo, time.Second*2)

		err := u.Store(context.TODO(), &mockArticle)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
	})
}
