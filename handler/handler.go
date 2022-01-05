package handler

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"math/rand"
	"shorter/storage"
	"time"
)

type Handler struct {
	DB *storage.DB
}

type handler struct {
	db *storage.DB
}

type Template struct {
	templates *template.Template
}

func NewHandler(h Handler) *handler {
	return &handler{
		db: h.DB,
	}
}

type Service interface {
	Register(h handler, g echo.Group)
}

func (h handler) REGISTER(group echo.Group, service Service) {
	service.Register(h, group)
}

func makeShort() string {
	rand.Seed(time.Now().Unix())
	var s string
	ra := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 6; i++ {
		s += string(ra[rand.Intn(len(ra))])
	}
	return s
}