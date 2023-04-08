package category

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"microservice/domain/entitie"
	"microservice/domain/handler/helper"
	_interface "microservice/domain/interface"
	"net/http"
	"strconv"
)

type categoryHandler struct {
	CategoryService _interface.CategoryService
	Logger          *zap.SugaredLogger
}

func ProvideCategoryHandler(logger *zap.SugaredLogger, service _interface.CategoryService) _interface.CategoryHandler {
	logger.Info("Executing ProvideCategoryHandler.")
	return &categoryHandler{
		CategoryService: service,
		Logger:          logger,
	}
}

func RegisterHooks(
	lifecycle fx.Lifecycle,
	handler _interface.CategoryHandler,
	echoInstance *echo.Echo,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				echoInstance.PUT("/category", handler.Update)
				echoInstance.POST("/category", handler.Store)
				echoInstance.GET("/category/:id", handler.GetByID)
				echoInstance.DELETE("/category/:id", handler.Delete)
				echoInstance.GET("/category/getAll", handler.GetAll)
				return nil
			},
			OnStop: func(context.Context) error {
				return nil
			},
		},
	)
}

func (a *categoryHandler) Delete(echoContext echo.Context) error {

	idP, err := strconv.Atoi(echoContext.Param("id"))
	if err != nil {
		return echoContext.JSON(http.StatusNotFound, _interface.ErrNotFound.Error())
	}

	id := uint(idP)
	ctx := echoContext.Request().Context()

	category, err := a.CategoryService.GetByID(ctx, id)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	err = a.CategoryService.Delete(ctx, category.ID)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	return echoContext.JSON(http.StatusOK, "")
}

func (a *categoryHandler) GetByID(echoContext echo.Context) error {

	idP, err := strconv.Atoi(echoContext.Param("id"))
	if err != nil {
		return echoContext.JSON(http.StatusNotFound, _interface.ErrNotFound.Error())
	}

	id := uint(idP)
	ctx := echoContext.Request().Context()

	art, err := a.CategoryService.GetByID(ctx, id)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	return echoContext.JSON(http.StatusOK, art)
}

func (a *categoryHandler) GetAll(echoContext echo.Context) error {

	ctx := echoContext.Request().Context()

	art, err := a.CategoryService.GetAll(ctx)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	return echoContext.JSON(http.StatusOK, art)
}

func (a *categoryHandler) Store(echoContext echo.Context) (err error) {
	var ent entitie.Category
	err = echoContext.Bind(&ent)
	if err != nil {
		return echoContext.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&ent); !ok {
		return echoContext.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := echoContext.Request().Context()
	err = a.CategoryService.Store(ctx, &ent)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	return echoContext.JSON(http.StatusCreated, ent)
}

func (a *categoryHandler) Update(echoContext echo.Context) (err error) {
	var ent entitie.Category
	err = echoContext.Bind(&ent)
	if err != nil {
		return echoContext.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&ent); !ok {
		return echoContext.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := echoContext.Request().Context()
	err = a.CategoryService.Update(ctx, &ent)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	return echoContext.JSON(http.StatusCreated, ent)
}

func isRequestValid(m *entitie.Category) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
