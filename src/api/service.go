package api

import (
	"net/http"

	"github.com/jat001/ddns4cdn/core"
	"github.com/labstack/echo/v4"
)

type service struct {
}

func (ctx *service) serviceStats(c echo.Context) error {
	return c.JSON(http.StatusOK, core.Store.ServiceStats2)
}

func (ctx *service) listServices(c echo.Context) error {
	s := map[string]bool{}
	core.Store.RunningService.Range(func(k, v any) bool {
		s[k.(string)] = v.(bool)
		return true
	})

	return c.JSON(http.StatusOK, s)
}

func Service(e *echo.Echo) {
	ctx := service{}

	g := e.Group("/service")

	g.GET("/", ctx.listServices)
	g.GET("/stats", ctx.serviceStats)
}
