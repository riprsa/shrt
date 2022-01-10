package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"shorter/config"
	"shorter/handler"
	"shorter/storage"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	t := &Template{
		templates: template.Must(template.ParseGlob("view/*.html")),
	}
	db, err := storage.Open()
	if err != nil {
		panic(err)
	}

	h := handler.NewHandler(handler.Handler{DB: db})

	e := echo.New()
	e.Renderer = t

	e.Use(middleware.Logger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	h.REGISTER(*e.Group(""), &handler.ShortsStorage{})

	c := config.Config{}
	err = c.LoadData()
	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":"+c.Port))
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
