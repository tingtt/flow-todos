package main

import (
	"flow-todos/jwt"
	"flow-todos/todo"
	"net/http"
	"strconv"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func complete(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*jwtIssuer, u)
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
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		return echo.ErrNotFound
	}
	if invalidUnit {
		// 409: Conflict
		return c.JSONPretty(http.StatusConflict, map[string]string{"message": "invalid todo repeat unit"}, "	")
	}
	if dateNotFound {
		// 409: Conflict
		return c.JSONPretty(http.StatusConflict, map[string]string{"message": "todo.date does not exists"}, "	")
	}

	if newTodo.Id != 0 {
		// 200: Success
		return c.JSONPretty(http.StatusOK, []todo.Todo{t, newTodo}, "	")
	}
	// 200: Success
	return c.JSONPretty(http.StatusOK, t, "	")
}