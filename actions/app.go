package actions

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	md "github.com/thejayhaykid/home-library/models"
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
func startPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

// SearchBook with a given search query
func searchBook(c echo.Context) error {
	book := md.Book
	return c.Render(http.StatusOK, "index", "World")
}

// PutBooks with ID
func putBooks(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

// GetBooks with a give filter
func getBooks(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

// DeleteBooks with given ID
func deleteBooks(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

// Login to account
func login(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

// Logout of account
func logout(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

// Start sets up the new instance of echo
func Start() *echo.Echo {
	t := &Template{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.GET("/", startPage)
	e.POST("/search", searchBook)
	e.PUT("/books", putBooks)
	e.GET("/books", getBooks)
	e.DELETE("/books/:pk", deleteBooks)
	e.GET("/login", login)
	e.GET("/logout", logout)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/assets", "assets")
	return e
}
