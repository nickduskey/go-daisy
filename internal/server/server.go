package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/nickduskey/go-daisy/internal/middleware"
	"github.com/nickduskey/go-daisy/internal/template"
)

func buildRoutes(e *echo.Echo) {
	// Handle static files
	e.File("/favicon.ico", "static/favicon.ico")
	e.Static("/static", "static")

	e.GET("/", HomeHandler)
}

func HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, template.Home("go-daisy"))
}

func BuildServerAddress() string {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	return fmt.Sprintf("%s:%s", host, port)
}

func Render(c echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	return c.HTML(statusCode, buf.String())
}

func Run() {
	e := echo.New()
	middleware.Apply(e)
	buildRoutes(e)
	address := BuildServerAddress()
	e.Logger.Fatal(e.Start(address))
}
