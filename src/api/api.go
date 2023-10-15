package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/jat001/ddns4cdn/core"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GraphQL(c echo.Context) error {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"history": historyQuery(),
			},
		}),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	req, _ := io.ReadAll(c.Request().Body)
	params := graphql.Params{
		Schema:        schema,
		RequestString: string(req),
	}

	res := graphql.Do(params)

	if len(res.Errors) > 0 {
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, res)

}

func API(config *core.Config) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/graphql", GraphQL)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.API.Port)))
}
