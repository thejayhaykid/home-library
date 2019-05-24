package actions

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// StartPage is a start
func StartPage(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// Start sets up the new instance of echo
func Start() *echo.Echo {
	e := echo.New()
	e.GET("/", StartPage)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return e
}
