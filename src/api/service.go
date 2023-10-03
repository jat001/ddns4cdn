package api

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

func listServices(ctx echo.Context, m *sync.Map) error {
	s := map[string]bool{}
	m.Range(func(k, v any) bool {
		s[k.(string)] = v.(bool)
		return true
	})

	return ctx.JSON(http.StatusOK, s)
}

func service(e *echo.Echo, m *sync.Map) {
	g := e.Group("/service")

	g.GET("/", func(ctx echo.Context) error {
		return listServices(ctx, m)
	})
}
