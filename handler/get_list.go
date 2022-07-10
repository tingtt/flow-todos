package handler

import (
	"flow-todos/flags"
	"flow-todos/jwt"
	"flow-todos/todo"
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type GetListQuery struct {
	Start               *string `query:"start" validate:"omitempty,datetime"`
	End                 *string `query:"end" validate:"omitempty,datetime"`
	ProjectId           *uint64 `query:"project_id" validate:"omitempty,gte=1"`
	WithCompleted       bool    `query:"with_completed" validate:"omitempty"`
	WithRepeatSchedules bool    `query:"with_repeat_schedules" validate:"omitempty"`
}

func datetimeStrConv(str string) (t time.Time, err error) {
	// y-m-dTh:m:s or unix timestamp
	t, err1 := time.Parse("2006-1-2T15:4:5", str)
	if err1 == nil {
		return
	}
	t, err2 := time.Parse(time.RFC3339, str)
	if err2 == nil {
		return
	}
	u, err3 := strconv.ParseInt(str, 10, 64)
	if err3 == nil {
		t = time.Unix(u, 0)
		return
	}
	err = fmt.Errorf("\"%s\" is not a unix timestamp or string format \"2006-1-2T15:4:5\"", str)
	return
}

func GetList(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
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
	if query.WithRepeatSchedules && query.End == nil {
		// 400: Bad request
		c.Logger().Debug("\"end\" required to get repeat schedules")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "\"end\" required to get repeat schedules"}, "	")
	}
	queryParsed := todo.GetListQuery{Start: start, End: end, ProjectId: query.ProjectId, WithCompleted: query.WithCompleted, WithRepeatSchedules: query.WithRepeatSchedules}

	// Get todos
	todos, err := todo.GetList(userId, queryParsed)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	if todos == nil {
		return c.JSONPretty(http.StatusOK, []interface{}{}, "	")
	}
	return c.JSONPretty(http.StatusOK, todos, "	")
}
