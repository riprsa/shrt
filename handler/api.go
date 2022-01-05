package handler

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
)

type APIShortsStorage struct {
	handler
}

func (a APIShortsStorage) Register(h handler, g echo.Group) {
	a.handler = h

	g.GET("", a.Wait)
	g.POST("", a.Create)
}

func (a APIShortsStorage) Create(c echo.Context) error {
	d, err := a.db.ByURL(c.FormValue("url"))

	if err == sql.ErrNoRows {
		for {
			ms := makeShort()
			if b, err := a.CheckCollision(ms); b {
				if err != nil {
					return c.JSON(http.StatusOK, err)
				}
				err = a.db.Insert(c.FormValue("url"), ms)
				if err != nil {
					return c.JSON(http.StatusOK, err)
				}
				return c.JSON(http.StatusOK, c.Request().Host+"/"+ms)
			}
		}
	} else if err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, c.Request().Host+"/"+d.Short)
}

func (a APIShortsStorage) CheckCollision(ms string) (bool, error) {
	_, err := a.db.ByShort(ms)
	if err == sql.ErrNoRows {
		return true, nil
	}
	return false, err
}

func (a APIShortsStorage) Wait(c echo.Context) error {
	return c.JSON(http.StatusOK, "to get something just make POST request with link with field url")
}
