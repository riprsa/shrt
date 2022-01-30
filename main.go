package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"shorter/internal/handler"
	"shorter/internal/html"
	"shorter/internal/service"
	"shorter/internal/storage"
)

func main() {
	db, err := storage.Open()
	if err != nil {
		panic(err)
	}

	s := service.New(db)

	h := handler.New(s)

	e := echo.New()
	t := html.New()
	e.Renderer = t
	e.Static("/assets", "view/assets")

	e.Use(middleware.Logger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	h.NewGroup(e.Group(""), &handler.Shorts{})
	h.NewGroup(e.Group("api"), &handler.ShortsAPI{})

	//e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

	pem := os.Getenv("SERVER_PEM")
	key := os.Getenv("SERVER_KEY")

	if err := e.StartTLS(":"+os.Getenv("PORT"), []byte(pem), []byte(key)); err != http.ErrServerClosed {
		panic(err)
	}
}
