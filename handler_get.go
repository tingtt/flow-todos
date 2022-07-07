package main

import (
	"flow-todos/jwt"
	"flow-todos/todo"
	"net/http"
	"strconv"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func get(c echo.Context) error {
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

	t, notFound, err := todo.Get(userId, id)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("todo not found")
		return echo.ErrNotFound
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, t, "	")
}
