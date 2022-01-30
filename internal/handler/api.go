package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"shorter/internal/model"
	"shorter/internal/validate"
)

type ShortsAPI struct {
	handler
}

func (s ShortsAPI) Register(h handler, g *echo.Group) {
	s.handler = h

	g.GET("", s.GetURL)
	g.POST("", s.ProcessURL)
}

func (s ShortsAPI) GetURL(c echo.Context) error {
	return c.JSON(http.StatusOK, "json POST request with field url")
}


func (s ShortsAPI) ProcessURL(c echo.Context) error {
	var data model.Data
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	url, ok := validate.URL(data.URL)
	if !ok {
		return c.JSON(http.StatusNotFound, "wrong_url")
	}

	short, err := s.service.GetShort(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, c.Request().Host+"/"+short)
}