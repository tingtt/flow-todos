package main

import (
	"flow-todos/flags"
	"flow-todos/handler"
	"flow-todos/jwt"
	"flow-todos/mysql"
	"flow-todos/todo"
	"flow-todos/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

type CustomValidator struct {
	validator *validator.Validate
}

func DatetimeStrValidation(fl validator.FieldLevel) bool {
	_, err1 := time.Parse("2006-1-2T15:4:5", fl.Field().String())
	_, err2 := time.Parse(time.RFC3339, fl.Field().String())
	_, err3 := strconv.ParseUint(fl.Field().String(), 10, 64)
	return err1 == nil || err2 == nil || err3 == nil
}

func (cv *CustomValidator) Validate(i interface{}) error {
	// Register custum validations
	cv.validator.RegisterValidation("datetime", DatetimeStrValidation)
	cv.validator.RegisterValidation("Y-M-D", todo.DateStrValidation)
	cv.validator.RegisterValidation("H:M", todo.HMTimeStrValidation)
	cv.validator.RegisterValidation("step15", todo.Step15IntValidation)

	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

func main() {
	// Get command line params / env variables
	f := flags.Get()

	//
	// Setup echo and middlewares
	//

	// Echo instance
	e := echo.New()

	// Gzip
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: int(*f.GzipLevel),
	}))

	// CORS
	if f.AllowOrigins != nil && len(f.AllowOrigins) != 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: f.AllowOrigins,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}))
	}

	// Log level
	e.Logger.SetLevel(log.Lvl(*f.LogLevel))

	// Validator instance
	e.Validator = &CustomValidator{validator: validator.New()}

	// JWT
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwt.JwtCustumClaims{},
		SigningKey: []byte(*f.JwtSecret),
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/-/readiness"
		},
	}))

	//
	// Setup DB
	//

	// DB client instance
	e.Logger.Info(mysql.SetDSNTCP(*f.MysqlUser, *f.MysqlPasswd, *f.MysqlHost, int(*f.MysqlPort), *f.MysqlDB))

	// Check connection
	d, err := mysql.Open()
	if err != nil {
		e.Logger.Fatal(err)
	}
	if err = d.Ping(); err != nil {
		e.Logger.Fatal(err)
	}

	//
	// Check health of external service
	//

	// flow-projects
	if *flags.Get().ServiceUrlProjects == "" {
		e.Logger.Fatal("`--service-url-projects` option is required")
	}
	if status, err := utils.HttpGet(*flags.Get().ServiceUrlProjects+"/-/readiness", nil); err != nil {
		e.Logger.Fatalf("failed to check health of external service `flow-projects` %s", err)
	} else if status != http.StatusOK {
		e.Logger.Fatal("failed to check health of external service `flow-projects`")
	}
	// flow-sprints
	if *flags.Get().ServiceUrlSprints == "" {
		e.Logger.Fatal("`--service-url-sprints` option is required")
	}
	if status, err := utils.HttpGet(*flags.Get().ServiceUrlSprints+"/-/readiness", nil); err != nil {
		e.Logger.Fatalf("failed to check health of external service `flow-sprints` %s", err)
	} else if status != http.StatusOK {
		e.Logger.Fatal("failed to check health of external service `flow-sprints`")
	}

	//
	// Routes
	//

	// Health check route
	e.GET("/-/readiness", func(c echo.Context) error {
		return c.String(http.StatusOK, "flow-todos is Healthy.\n")
	})

	// Restricted routes
	e.GET("/", handler.GetList)
	e.POST("/", handler.Post)
	e.GET(":id", handler.Get)
	e.PATCH(":id", handler.Patch)
	e.DELETE(":id", handler.Delete)
	e.PATCH(":id/skip", handler.Skip)
	e.PATCH(":id/complete", handler.Complete)
	e.DELETE("/", handler.DeleteAll)

	//
	// Start echo
	//
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *f.Port)))
}
