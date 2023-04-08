package tag

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	customEchoContext "microservice/cmd/infra/context"
	"microservice/domain/entitie"
	"microservice/domain/handler/helper"
	_interface "microservice/domain/interface"
	"net/http"
	"strconv"
)

type tagHandler struct {
	TagService _interface.TagService
	Logger     *zap.SugaredLogger
}

func ProvideTagHandler(logger *zap.SugaredLogger, service _interface.TagService) _interface.TagHandler {
	logger.Info("Executing ProvideTagHandler.")
	return &tagHandler{
		TagService: service,
		Logger:     logger,
	}
}

func RegisterHooks(
	lifecycle fx.Lifecycle,
	handler _interface.TagHandler,
	echoInstance *echo.Echo,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				echoInstance.PUT("/tag", handler.Update)
				echoInstance.POST("/tag", handler.Store)
				echoInstance.GET("/tag/:id", handler.GetByID)
				echoInstance.DELETE("/tag/:id", handler.Delete)
				echoInstance.GET("/tag/getAll", handler.GetAll)
				return nil
			},
			OnStop: func(context.Context) error {
				return nil
			},
		},
	)
}

func (a tagHandler) Delete(echoContext echo.Context) error {

	idP, err := strconv.Atoi(echoContext.Param("id"))
	if err != nil {
		return echoContext.JSON(http.StatusNotFound, _interface.ErrNotFound.Error())
	}

	id := uint(idP)
	ctx := echoContext.Request().Context()

	tag, err := a.TagService.GetByID(ctx, id)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	err = a.TagService.Delete(ctx, tag.ID)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	return echoContext.JSON(http.StatusOK, "")
}

func (a tagHandler) GetByID(echoContext echo.Context) error {

	idP, err := strconv.Atoi(echoContext.Param("id"))
	if err != nil {
		return echoContext.JSON(http.StatusNotFound, _interface.ErrNotFound.Error())
	}

	id := uint(idP)
	ctx := echoContext.Request().Context()

	art, err := a.TagService.GetByID(ctx, id)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	return echoContext.JSON(http.StatusOK, art)
}

func (a tagHandler) GetAll(echoContext echo.Context) error {

	cc := echoContext.(*customEchoContext.EchoContext)
	a.Logger.Info("Correlation", cc.GetCorrelationID())

	ctx := echoContext.Request().Context()

	art, err := a.TagService.GetAll(ctx)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	return echoContext.JSON(http.StatusOK, art)
}

func (a tagHandler) Store(echoContext echo.Context) (err error) {
	var ent entitie.Tag
	err = echoContext.Bind(&ent)
	if err != nil {
		return echoContext.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&ent); !ok {
		return echoContext.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := echoContext.Request().Context()
	err = a.TagService.Store(ctx, &ent)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	return echoContext.JSON(http.StatusCreated, ent)
}

func (a tagHandler) Update(echoContext echo.Context) (err error) {
	var ent entitie.Tag
	err = echoContext.Bind(&ent)
	if err != nil {
		return echoContext.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&ent); !ok {
		return echoContext.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := echoContext.Request().Context()
	err = a.TagService.Update(ctx, &ent)
	if err != nil {
		return echoContext.JSON(helper.GetStatusCode(err), helper.ResponseError{Message: err.Error()})
	}

	return echoContext.JSON(http.StatusCreated, ent)
}

func isRequestValid(m *entitie.Tag) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
