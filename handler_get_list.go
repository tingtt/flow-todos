package main

import (
	"flow-todos/jwt"
	"flow-todos/todo"
	"net/http"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type GetListQuery struct {
	Start         *string `query:"start" validate:"omitempty,datetime"`
	End           *string `query:"end" validate:"omitempty,datetime"`
	ProjectId     *uint64 `query:"project_id" validate:"omitempty,gte=1"`
	WithCompleted bool    `query:"with_completed" validate:"omitempty"`
}

func getList(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*jwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// Bind query
	query := new(GetListQuery)
	if err = c.Bind(query); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate query
	if err = c.Validate(query); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}
	var start, end *time.Time
	if query.Start != nil {
		startTmp, err := datetimeStrConv(*query.Start)
		if err != nil {
			// 400: Bad request
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}
		start = &startTmp
	}
	if query.End != nil {
		endTmp, err := datetimeStrConv(*query.End)
		if err != nil {
			// 400: Bad request
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}
		end = &endTmp
	}
	queryParsed := todo.GetListQuery{Start: start, End: end, ProjectId: query.ProjectId, WithCompleted: query.WithCompleted}

	// Get todos
	todos, err := todo.GetList(userId, queryParsed)
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	if todos == nil {
		return c.JSONPretty(http.StatusOK, []interface{}{}, "	")
	}
	return c.JSONPretty(http.StatusOK, todos, "	")
}
