package main

import (
	"fmt"
	"net/http"
	"os"

	"shorter/internal/handler"
	"shorter/internal/service"
	"shorter/internal/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	mode := ""
	fmt.Println("choose mode:\n1. without pem\n2. with pem")
	fmt.Scan(&mode)

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

	if mode == "1" {
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
