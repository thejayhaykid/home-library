package actions

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Template is for rendering
type Template struct {
	templates *template.Template
}

// Render render's the template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// StartPage is a start
func StartPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

// Start sets up the new instance of echo
func Start() *echo.Echo {
	t := &Template{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.GET("/", StartPage)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/assets", "assets")
	return e
}
