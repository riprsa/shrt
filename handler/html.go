package handler

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type ShortsStorage struct {
	handler
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
		ms := makeShort()
		for {
			if b, err := s.CheckCollision(ms); b {
				if err != nil {
					return c.Render(http.StatusOK, "error", err)
				}
				err = s.db.Insert(c.FormValue("url"), ms)
				if err != nil {
					return c.Render(http.StatusOK, "error", err)
				}
				return c.Render(http.StatusOK, "result", c.Request().Host+"/"+ms)
			}
		}
	} else if err != nil {
		return c.Render(http.StatusOK, "error", err)
	}
	return c.Render(http.StatusOK, "result", c.Request().Host+"/"+d.Short)
}

func (s ShortsStorage) CheckCollision(ms string) (bool, error) {
	_, err := s.db.ByShort(ms)
	if err == sql.ErrNoRows {
		return true, nil
	}
	return false, err
}

func (s ShortsStorage) Redirect(c echo.Context) error {
	data, err := s.db.ByShort(c.Request().URL.Path[1:])

	if err != nil && data.URL != "" {
		return c.Redirect(http.StatusFound, "/")
	}

	data.URL = strings.TrimPrefix(data.URL, "https://")
	data.URL = strings.TrimPrefix(data.URL, "http://")
	data.URL = strings.TrimPrefix(data.URL, "ftp://")

	return c.Redirect(http.StatusFound, "https://"+data.URL)
}
