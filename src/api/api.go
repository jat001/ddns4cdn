package api

import (
	"fmt"

	"github.com/jat001/ddns4cdn/core"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func API(config *core.Config) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	GraphQL(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.API.Port)))
}
