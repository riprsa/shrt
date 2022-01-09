package handler

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
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
	url, ok := s.val.URLValidation(c.FormValue("url"))
	if !ok {
		return c.Render(http.StatusNotFound, "wrongUrl", nil)
	}

	dataFromDB, err := s.db.ByURL(url)
	if err == sql.ErrNoRows {
		short, err := s.CreateShort(url)
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "result", c.Request().Host+"/"+short)
	} else if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "result", c.Request().Host+"/"+dataFromDB.Short)
}

func (s ShortsStorage) CreateShort(url string) (string, error) {
	for {
		short := makeShort()
		if b, err := s.CheckCollision(short); b {
			if err != nil {
				return "", err
			}
			err = s.db.Insert(url, short)
			if err != nil {
				return "", err
			}
			return short, nil
		}
	}
}

func (s ShortsStorage) CheckCollision(ms string) (bool, error) {
	_, err := s.db.ByShort(ms)
	if err == sql.ErrNoRows {
		return true, nil
	}
	return false, err
}

func (s ShortsStorage) Redirect(c echo.Context) error {
	short := c.Request().URL.Path[1:]
	data, err := s.db.ByShort(short)
	if err != nil {
		return err
	}

	url, ok := s.val.URLValidation(data.URL)
	if !ok {
		return c.Render(http.StatusNotFound, "404", nil)
	}
	return c.Redirect(http.StatusFound, url)
}
