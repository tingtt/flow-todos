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

func Skip(c echo.Context) error {
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

	t, overUntil, notFound, repeatNotFound, dateNotFound, invalidUnit, err := todo.Skip(userId, id)
	if err != nil {
		// 500: Internal Server Error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		return echo.ErrNotFound
	}
	if repeatNotFound {
		// 400: Bad request
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "repeat not found"}, "	")
	}
	if invalidUnit {
		// 400: Bad request
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "invalid todo repeat unit"}, "	")
	}
	if dateNotFound {
		// 400: Bad request
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "todo.date does not exists"}, "	")
	}
	if overUntil {
		// 400: Bad request
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "cannot skip last todo in due date"}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, t, "	")
}
