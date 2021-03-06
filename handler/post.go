package handler

import (
	"flow-todos/flags"
	"flow-todos/jwt"
	"flow-todos/todo"
	"flow-todos/utils"
	"fmt"
	"net/http"
	"strings"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Post(c echo.Context) error {
	// Check `Content-Type`
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		// 415: Invalid `Content-Type`
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": "unsupported media type"}, "	")
	}

	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// Bind request body
	post := new(todo.PostBody)
	if err = c.Bind(post); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate request body
	if err = c.Validate(post); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// Check project id
	if post.ProjectId != nil {
		status, err := utils.HttpGet(fmt.Sprintf("%s/%d", *flags.Get().ServiceUrlProjects, *post.ProjectId), &u.Raw)
		if err != nil {
			// 500: Internal server error
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if status != http.StatusOK {
			// 400: Bad request
			c.Logger().Debug(fmt.Sprintf("project id: %d does not exist", *post.ProjectId))
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("project id: %d does not exist", *post.ProjectId)}, "	")
		}
	}

	// Check sprint id
	if post.SprintId != nil {
		status, err := utils.HttpGet(fmt.Sprintf("%s/%d", *flags.Get().ServiceUrlSprints, *post.SprintId), &u.Raw)
		if err != nil {
			// 500: Internal server error
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if status != http.StatusOK {
			// 400: Bad request
			c.Logger().Debugf("sprint id: %d does not exist", *post.SprintId)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("sprint id: %d does not exist", *post.SprintId)}, "	")
		}
	}

	p, dateNotFound, dateOverUntil, noDaysWithWeekly, err := todo.Post(userId, *post)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if dateNotFound {
		// 400: Bad request
		c.Logger().Debug("`date` required to set `repeat.until`")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "`date` required to set `repeat`"}, "	")
	}
	if dateOverUntil {
		// 400: Bad request
		c.Logger().Debug("`date` must until `repeat.until`")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "`date` must until `repeat.until`"}, "	")
	}
	if noDaysWithWeekly {
		// 400: Bad request
		c.Logger().Debug("`repeat.days` required with `repeat.unit: \"week\"`")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "`repeat.days` required with `repeat.unit: \"week\"`"}, "	")
	}

	// 201: Created
	return c.JSONPretty(http.StatusCreated, p, "	")
}
