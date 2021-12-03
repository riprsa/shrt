package handler

import (
	"database/sql"
	"html/template"
	"math/rand"
	"net/http"
	"shorter/storage"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
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

type ShortsStorage struct {
	handler
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

func (s ShortsStorage) Register(h handler, g echo.Group) {
	s.handler = h

	g.GET("/", s.Wait)
	g.POST("/", s.Create)
	g.GET("*", s.Redirect)
}

func (s ShortsStorage) Wait(c echo.Context) error {
	return c.Render(http.StatusOK, "main", nil)
}

func (s ShortsStorage) Create(c echo.Context) error {
	d, err := s.db.ByURL(c.FormValue("url"))
	if err == sql.ErrNoRows {
		make: ms := makeShort()
		if b, err := s.CheckCollision(ms); b {
			if err != nil {
				return err
			}
			err = s.db.Insert(c.FormValue("url"), ms)
			if err != nil {
				return c.Render(http.StatusOK, "error", err)
			}
			return c.Render(http.StatusOK, "result", c.Request().Host+"/"+ms)
		} else {
			goto make
		}
	} else if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "result", c.Request().Host+"/"+d.Short)
}

func (s ShortsStorage) CheckCollision(ms string) (bool, error) {
	_, err := s.db.ByShort(ms)
	if err == sql.ErrNoRows {
		return true, nil
	} else if err != nil {
		return false, err
	}
	return false, nil
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

func (s ShortsStorage) Redirect(c echo.Context) error {
	data, err := s.db.ByShort(c.Request().URL.Path[1:])
	if err != nil && data.URL != ""{
		return c.Redirect(http.StatusFound, "/")
	}

	data.URL = strings.TrimPrefix(data.URL, "https://")
	data.URL = strings.TrimPrefix(data.URL, "http://")
	data.URL = strings.TrimPrefix(data.URL, "ftp://")

	err = c.Redirect(http.StatusFound, "https://"+data.URL)
	return err
}