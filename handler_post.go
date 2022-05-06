package main

import (
	"flow-todos/jwt"
	"flow-todos/todo"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func post(c echo.Context) error {
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
	if post.Repeat != nil && post.Repeat.Until != nil {
		// Validate `date` and `repeat.until`
		if post.Date == nil {
			// 400: Bad request
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "`date` required to set `repeat.until`"}, "	")
		}
		var t1, t2 time.Time
		t1, _ = time.Parse("2006-1-2", *post.Date)
		t2, _ = time.Parse("2006-1-2", *post.Repeat.Until)
		if t1.After(t2) {
			// 400: Bad request
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "`date` must until `repeat.until`"}, "	")
		}
	}

	// Check project id
	if post.ProjectId != nil {
		valid, err := checkProjectId(u.Raw, *post.ProjectId)
		if err != nil {
			// 500: Internal server error
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if !valid {
			// 409: Conflit
			c.Logger().Debug(fmt.Sprintf("project id: %d does not exist", *post.ProjectId))
			return c.JSONPretty(http.StatusConflict, map[string]string{"message": fmt.Sprintf("project id: %d does not exist", *post.ProjectId)}, "	")
		}
	}

	// Check sprint id
	if post.SprintId != nil {
		valid, err := checkSprintId(u.Raw, *post.SprintId)
		if err != nil {
			// 500: Internal server error
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if !valid {
			// 409: Conflit
			c.Logger().Debug(fmt.Sprintf("sprint id: %d does not exist", *post.SprintId))
			return c.JSONPretty(http.StatusConflict, map[string]string{"message": fmt.Sprintf("sprint id: %d does not exist", *post.SprintId)}, "	")
		}
	}

	p, err := todo.Post(userId, *post)
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 201: Created
	return c.JSONPretty(http.StatusCreated, p, "	")
}
