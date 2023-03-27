package handler

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	_interface "microservice/domain/interface"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	zap.L().Error("error", zap.Error(err))

	switch err {
	case gorm.ErrRecordNotFound:
		return http.StatusNotFound
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
