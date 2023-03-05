package context

import "github.com/labstack/echo/v4"

type EchoContext struct {
	echo.Context
}

func (c *EchoContext) Foo() {
	println("foo")
}

func (c *EchoContext) Bar() {
	println("bar")
}
