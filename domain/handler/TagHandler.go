package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	_interface "microservice/domain/interface"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

type TagHandler struct {
	TagService _interface.TagService
}

func NewTagHandler(e *echo.Echo, service _interface.TagService) {
	handler := &TagHandler{
		TagService: service,
	}
	e.GET("/tag/:id", handler.GetByID)
}

func (a *TagHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, _interface.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	art, err := a.TagService.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	zap.L().Error("error", zap.Error(err))

	switch err {
	case _interface.ErrInternalServerError:
		return http.StatusInternalServerError
	case _interface.ErrNotFound:
		return http.StatusNotFound
	case _interface.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
