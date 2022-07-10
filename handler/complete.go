package handler

import (
	"flow-todos/flags"
	"flow-todos/jwt"
	"flow-todos/todo"
	"net/http"
	"strconv"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Complete(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// id
	idStr := c.Param("id")

	// string -> uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 404: Not found
		return echo.ErrNotFound
	}

	t, newTodo, notFound, dateNotFound, invalidUnit, err := todo.Complete(userId, id)
	if err != nil {
		// 500: Internal Server Error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		return echo.ErrNotFound
	}
	if invalidUnit {
		// 400: Bad request
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "invalid todo repeat unit"}, "	")
	}
	if dateNotFound {
		// 400: Bad request
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "todo.date does not exists"}, "	")
	}

	if newTodo.Id != 0 {
		// 200: Success
		return c.JSONPretty(http.StatusOK, []todo.Todo{t, newTodo}, "	")
	}
	// 200: Success
	return c.JSONPretty(http.StatusOK, t, "	")
}
