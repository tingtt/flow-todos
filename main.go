package main

import (
	"flag"
	"flow-todos/jwt"
	"flow-todos/mysql"
	"flow-todos/todo"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func getIntEnv(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		var intValue, err = strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Priority: command line params > env variables > default value
var (
	port               = flag.Int("port", getIntEnv("PORT", 1323), "Server port")
	logLevel           = flag.Int("log-level", getIntEnv("LOG_LEVEL", 2), "Log level (1: 'DEBUG', 2: 'INFO', 3: 'WARN', 4: 'ERROR', 5: 'OFF', 6: 'PANIC', 7: 'FATAL'")
	gzipLevel          = flag.Int("gzip-level", getIntEnv("GZIP_LEVEL", 6), "Gzip compression level")
	mysqlHost          = flag.String("mysql-host", getEnv("MYSQL_HOST", "db"), "MySQL host")
	mysqlPort          = flag.Int("mysql-port", getIntEnv("MYSQL_PORT", 3306), "MySQL port")
	mysqlDB            = flag.String("mysql-database", getEnv("MYSQL_DATABASE", "flow-todos"), "MySQL database")
	mysqlUser          = flag.String("mysql-user", getEnv("MYSQL_USER", "flow-todos"), "MySQL user")
	mysqlPasswd        = flag.String("mysql-password", getEnv("MYSQL_PASSWORD", ""), "MySQL password")
	jwtIssuer          = flag.String("jwt-issuer", getEnv("JWT_ISSUER", "flow-users"), "JWT issuer")
	jwtSecret          = flag.String("jwt-secret", getEnv("JWT_SECRET", ""), "JWT secret")
	serviceUrlProjects = flag.String("service-url-projects", getEnv("SERVICE_URL_PROJECTS", ""), "Service url: flow-projects")
	serviceUrlSprints  = flag.String("service-url-sprints", getEnv("SERVICE_URL_SPRINTS", ""), "Service url: flow-sprints")
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
	flag.Parse()
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: *gzipLevel,
	}))
	e.Logger.SetLevel(log.Lvl(*logLevel))
	e.Validator = &CustomValidator{validator: validator.New()}

	// Setup db client instance
	e.Logger.Info(mysql.SetDSNTCP(*mysqlUser, *mysqlPasswd, *mysqlHost, *mysqlPort, *mysqlDB))
	// Check connection
	d, err := mysql.Open()
	if err != nil {
		e.Logger.Fatal(err)
	}
	if err = d.Ping(); err != nil {
		e.Logger.Fatal(err)
	}

	// Service status check
	if *serviceUrlProjects == "" {
		e.Logger.Fatal("`--service-url-projects` option is required")
	}
	if ok, err := checkHealth(*serviceUrlProjects + "/-/readiness"); err != nil {
		e.Logger.Fatalf("failed to check health of external service `flow-projects` %s", err)
	} else if !ok {
		e.Logger.Fatal("failed to check health of external service `flow-projects`")
	}
	if *serviceUrlSprints == "" {
		e.Logger.Fatal("`--service-url-sprints` option is required")
	}
	if ok, err := checkHealth(*serviceUrlSprints + "/-/readiness"); err != nil {
		e.Logger.Fatalf("failed to check health of external service `flow-sprints` %s", err)
	} else if !ok {
		e.Logger.Fatal("failed to check health of external service `flow-sprints`")
	}

	// Setup JWT
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwt.JwtCustumClaims{},
		SigningKey: []byte(*jwtSecret),
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/-/readiness"
		},
	}))

	// Health check route
	e.GET("/-/readiness", func(c echo.Context) error {
		return c.String(http.StatusOK, "flow-todos is Healthy.\n")
	})

	// Restricted routes
	e.GET("/", getList)
	e.POST("/", post)
	e.GET(":id", get)
	e.PATCH(":id", patch)
	e.DELETE(":id", delete)
	e.PATCH(":id/skip", skip)
	e.PATCH(":id/complete", complete)
	e.DELETE("/", deleteAll)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
