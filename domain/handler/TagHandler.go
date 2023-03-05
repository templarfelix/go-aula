package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
	"microservice/domain/entitie"
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
	e.POST("/tag", handler.Store)
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

func (a *TagHandler) Store(c echo.Context) (err error) {
	var ent entitie.Tag
	err = c.Bind(&ent)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&ent); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.TagService.Store(ctx, &ent)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, ent)
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

func isRequestValid(m *entitie.Tag) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
