package main

import (
	"flow-todos/jwt"
	"flow-todos/todo"
	"fmt"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func post(c echo.Context) error {
	// Check `Content-Type`
	if c.Request().Header.Get("Content-Type") != "application/json" &&
		c.Request().Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
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
	if post.TermId != nil {
		valid, err := checkProjectId(u.Raw, *post.TermId)
		if err != nil {
			// 500: Internal server error
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if !valid {
			// 409: Conflit
			c.Logger().Debug(fmt.Sprintf("project id: %d does not exist", *post.TermId))
			return c.JSONPretty(http.StatusConflict, map[string]string{"message": fmt.Sprintf("project id: %d does not exist", *post.TermId)}, "	")
		}
	}

	p, err := todo.Post(userId, *post)
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, p, "	")
}
