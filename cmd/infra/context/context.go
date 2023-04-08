package context

import "github.com/labstack/echo/v4"

type EchoContext struct {
	echo.Context
}

func (echoContext *EchoContext) GetCorrelationID() string {
	return echoContext.Request().Header.Get("x-correlation-id")
}
