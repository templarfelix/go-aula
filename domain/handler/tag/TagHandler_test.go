package tag_test

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"microservice/domain/entitie"
	"microservice/domain/handler/tag"
	"microservice/domain/interface/mocks"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
)

func TestGetByID(t *testing.T) {
	var mockArticle entitie.Tag
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)

	mockUCase := new(mocks.TagService)

	num := int(mockArticle.ID)

	mockUCase.On("GetByID", mock.Anything, uint(num)).Return(mockArticle, nil)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/tag/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))

	zap := zap.NewExample()
	slogger := zap.Sugar()

	handler := tag.ProvideTagHandler(slogger, mockUCase)

	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}
