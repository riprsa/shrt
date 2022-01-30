package handler

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"shorter/internal/model"
	"shorter/internal/validate"
)

type Shorts struct {
	handler
}

func (s Shorts) Register(h handler, g *echo.Group) {
	s.handler = h

	g.GET("/", s.Wait)
	g.POST("/", s.Create)

	g.GET("/api", s.WaitAPI)
	g.POST("/api", s.CreateAPI)

	g.GET("*", s.Redirect)
}

func (s Shorts) Wait(c echo.Context) error {
	return c.Render(http.StatusOK, "main", nil)
}

func (s Shorts) WaitAPI(c echo.Context) error {
	return c.JSON(http.StatusOK, "json POST request with field url")
}

func (s Shorts) Create(c echo.Context) error {
	url, ok := validate.URL(c.FormValue("url"))
	if !ok {
		return c.Render(http.StatusNotFound, "wrongURL", nil)
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

func (s Shorts) CreateAPI(c echo.Context) error {
	var data model.Data
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	url, ok := validate.URL(data.URL)
	if !ok {
		return c.JSON(http.StatusNotFound, "wrong_URL")
	}

	dataFromDB, err := s.db.ByURL(url)
	if err == sql.ErrNoRows {
		short, err := s.CreateShort(url)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, c.Request().Host+"/"+short)
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, c.Request().Host+"/"+dataFromDB.Short)
}

func (s Shorts) CreateShort(url string) (string, error) {
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

func (s Shorts) CheckCollision(ms string) (bool, error) {
	_, err := s.db.ByShort(ms)
	if err == sql.ErrNoRows {
		return true, nil
	}
	return false, err
}

func (s Shorts) Redirect(c echo.Context) error {
	short := c.Request().URL.Path[1:]
	data, err := s.db.ByShort(short)
	if err != nil {
		return err
	}

	url, ok := validate.URL(data.URL)
	if !ok {
		return c.Render(http.StatusNotFound, "404", nil)
	}
	return c.Redirect(http.StatusFound, url)
}
