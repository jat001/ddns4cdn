package api

import (
	"fmt"
	"sync"

	"github.com/jat001/ddns4cdn/core"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func API(config *core.Config, m *sync.Map) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	service(e, m)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.API.Port)))
}
