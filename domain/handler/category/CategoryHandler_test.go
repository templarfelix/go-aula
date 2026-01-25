package category_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"microservice/domain/entitie"
	"microservice/domain/handler/category"
	_interface "microservice/domain/interface"
	"microservice/domain/interface/mocks"
)

func setupTest() (*echo.Echo, *zap.SugaredLogger) {
	e := echo.New()
	logger := zap.NewExample().Sugar()
	return e, logger
}

func TestGetByID(t *testing.T) {
	var mockCategory entitie.Category
	err := faker.FakeData(&mockCategory)
	assert.NoError(t, err)

	e, logger := setupTest()
	mockService := new(mocks.CategoryService)

	t.Run("success", func(t *testing.T) {
		num := int(mockCategory.ID)
		mockService.On("GetByID", mock.Anything, uint(num)).Return(mockCategory, nil).Once()

		req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/category/"+strconv.Itoa(num), nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(num))

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid_id", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/category/invalid", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("not_found", func(t *testing.T) {
		mockService.On("GetByID", mock.Anything, uint(999)).Return(entitie.Category{}, _interface.ErrNotFound).Once()

		req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/category/999", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAll(t *testing.T) {
	e, logger := setupTest()
	mockService := new(mocks.CategoryService)

	t.Run("success", func(t *testing.T) {
		mockCategories := []entitie.Category{
			{Name: "Category 1"},
			{Name: "Category 2"},
		}
		mockService.On("GetAll", mock.Anything).Return(mockCategories, nil).Once()

		req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/category/getAll", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category/getAll")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.GetAll(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockService.On("GetAll", mock.Anything).Return([]entitie.Category{}, errors.New("database error")).Once()

		req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/category/getAll", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category/getAll")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.GetAll(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	e, logger := setupTest()
	mockService := new(mocks.CategoryService)

	t.Run("success", func(t *testing.T) {
		mockService.On("Store", mock.Anything, mock.AnythingOfType("*entitie.Category")).Return(nil).Once()

		req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/category",
			strings.NewReader(`{"name":"Test Category"}`))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.Store(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid_json", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/category",
			strings.NewReader(`invalid json`))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.Store(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("validation_error", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/category",
			strings.NewReader(`{"name":""}`))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.Store(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		mockService.On("Store", mock.Anything, mock.AnythingOfType("*entitie.Category")).
			Return(_interface.ErrConflict).Once()

		req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/category",
			strings.NewReader(`{"name":"Existing Category"}`))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.Store(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	e, logger := setupTest()
	mockService := new(mocks.CategoryService)

	t.Run("success", func(t *testing.T) {
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*entitie.Category")).Return(nil).Once()

		req, err := http.NewRequestWithContext(context.TODO(), echo.PUT, "/category",
			strings.NewReader(`{"id":1,"name":"Updated Category"}`))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid_json", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.TODO(), echo.PUT, "/category",
			strings.NewReader(`invalid json`))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("validation_error", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.TODO(), echo.PUT, "/category",
			strings.NewReader(`{"id":1,"name":""}`))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestDelete(t *testing.T) {
	e, logger := setupTest()
	mockService := new(mocks.CategoryService)

	t.Run("success", func(t *testing.T) {
		mockCategory := entitie.Category{Name: "Test"}
		mockCategory.ID = 1
		mockService.On("GetByID", mock.Anything, uint(1)).Return(mockCategory, nil).Once()
		mockService.On("Delete", mock.Anything, uint(1)).Return(nil).Once()

		req, err := http.NewRequestWithContext(context.TODO(), echo.DELETE, "/category/1", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid_id", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.TODO(), echo.DELETE, "/category/invalid", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("not_found", func(t *testing.T) {
		mockService.On("GetByID", mock.Anything, uint(999)).Return(entitie.Category{}, _interface.ErrNotFound).Once()

		req, err := http.NewRequestWithContext(context.TODO(), echo.DELETE, "/category/999", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/category/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")

		handler := category.ProvideCategoryHandler(logger, mockService)
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})
}
