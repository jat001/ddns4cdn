package api

import (
	"io"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/jat001/ddns4cdn/core"
	"github.com/labstack/echo/v4"
)

func query(c echo.Context) error {
	fields := graphql.Fields{
		"history": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "service",
				Fields: graphql.Fields{
					"id": &graphql.Field{
						Type: graphql.String,
					},
					"type": &graphql.Field{
						Type: graphql.String,
					},
					"history": &graphql.Field{
						Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
							Name: "history",
							Fields: graphql.Fields{
								"success": &graphql.Field{
									Type: graphql.Boolean,
								},
								"end_time": &graphql.Field{
									Type: graphql.DateTime,
								},
							},
						})),
					},
				},
			})),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"type": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				id, ok := p.Args["id"].(string)
				if !ok {
					id = ""
				}
				typ, ok := p.Args["type"].(string)
				if !ok {
					typ = ""
				}

				x := make([]map[string]any, 0, len(core.Store.ServiceMap))

				for _, i := range core.Store.ServiceMap {
					if id != "" && i[0].ID != id {
						continue
					}
					if typ != "" && i[0].Type != typ {
						continue
					}

					y := make([]map[string]any, 0, len(i))

					for _, j := range i {
						y = append(y, map[string]any{
							"success":  j.Success,
							"end_time": j.EndTime,
						})
					}

					x = append(x, map[string]any{
						"id":      i[0].ID,
						"type":    i[0].Type,
						"history": y,
					})
				}
				return x, nil
			},
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: fields,
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

func GraphQL(e *echo.Echo) {
	e.POST("/graphql", query)
}
