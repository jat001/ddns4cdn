package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/jat001/ddns4cdn/core"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GraphQLError(message string) map[string][]map[string]string {
	return map[string][]map[string]string{
		"errors": {
			{
				"message": message,
			},
		},
	}
}

func GraphQL(c echo.Context) error {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "query",
			Fields: graphql.Fields{
				"history": historyQuery(),
			},
		}),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GraphQLError(err.Error()))
	}

	req := make(map[string]any)
	err = json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GraphQLError("Problems parsing JSON: "+err.Error()))
	}

	query, ok := req["query"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, GraphQLError("A query attribute must be specified and must be a string."))
	}

	params := graphql.Params{
		Schema:        schema,
		RequestString: query,
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
	e.Use(middleware.CORS())

	e.POST("/graphql", GraphQL)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.API.Port)))
}
