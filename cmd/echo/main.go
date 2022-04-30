package main

import (
	"net/http"
	"os"

	handler "shorter/internal/handler/echo"
	"shorter/internal/service"
	"shorter/internal/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := storage.Open()
	if err != nil {
		panic(err)
	}

	s := service.New(db)

	h := handler.New(s)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	h.NewGroup(e.Group("api"), &handler.ShortService{})

	if os.Getenv("MODE") == "1" {
		e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
	} else {
		pem := os.Getenv("SERVER_PEM")
		key := os.Getenv("SERVER_KEY")

		if pem == "" || key == "" {
			panic("pem or key is empty")
		}

		if err := e.StartTLS(":"+os.Getenv("PORT"), []byte(pem), []byte(key)); err != http.ErrServerClosed {
			panic(err)
		}
	}
}
