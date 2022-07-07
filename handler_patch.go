package main

import (
	"flow-todos/jwt"
	"flow-todos/todo"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func patch(c echo.Context) error {
	// Check `Content-Type`
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		// 415: Invalid `Content-Type`
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": "unsupported media type"}, "	")
	}

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

	// Bind request body
	patch := new(todo.PatchBody)
	if err = c.Bind(patch); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate request body
	if err = c.Validate(patch); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// Check project id
	if patch.ProjectId.UInt64 != nil && *patch.ProjectId.UInt64 != nil {
		valid, err := checkProjectId(u.Raw, **patch.ProjectId.UInt64)
		if err != nil {
			// 500: Internal server error
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if !valid {
			// 409: Conflit
			c.Logger().Debugf("project id: %d does not exist", **patch.ProjectId.UInt64)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("project id: %d does not exist", **patch.ProjectId.UInt64)}, "	")
		}
	}

	// Check sprint id
	if patch.SprintId.UInt64 != nil && *patch.SprintId.UInt64 != nil {
		valid, err := checkSprintId(u.Raw, **patch.SprintId.UInt64)
		if err != nil {
			// 500: Internal server error
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if !valid {
			// 409: Conflit
			c.Logger().Debug(fmt.Sprintf("sprint id: %d does not exist", **patch.SprintId.UInt64))
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("sprint id: %d does not exist", **patch.SprintId.UInt64)}, "	")
		}
	}

	p, notFound, dateNotFound, dateOverUtil, noDaysWithWeekly, err := todo.Patch(userId, id, *patch)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("project not found")
		return echo.ErrNotFound
	}
	if dateNotFound {
		// 400: Bad request
		c.Logger().Debug("`date` required to set `repeat.until`")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "`date` required to set `repeat`"}, "	")
	}
	if dateOverUtil {
		// 400: Bad request
		c.Logger().Debug("`date` must until `repeat.until`")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "`date` must until `repeat.until`"}, "	")
	}
	if noDaysWithWeekly {
		// 400: Bad request
		c.Logger().Debug("`repeat.days` required with `repeat.unit: \"week\"`")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "`repeat.days` required with `repeat.unit: \"week\"`"}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, p, "	")
}
