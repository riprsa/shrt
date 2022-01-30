package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"shorter/internal/validate"
)

type Shorts struct {
	handler
}

func (s Shorts) Register(h handler, g *echo.Group) {
	s.handler = h

	g.GET("/", s.GetURL)
	g.POST("/", s.ProcessURL)

	g.GET("*", s.Redirect)
}

func (s Shorts) GetURL(c echo.Context) error {
	return c.Render(http.StatusOK, "main", nil)
}

func (s Shorts) ProcessURL(c echo.Context) error {
	url, ok := validate.URL(c.FormValue("url"))
	if !ok {
		return c.Render(http.StatusBadRequest, "wrongURL", nil)
	}

	short, err := s.service.GetShort(url)
	if err != nil {
		return c.Render(http.StatusBadRequest, "error", err)
	}

	return c.Render(http.StatusOK, "result", c.Request().Host+"/"+short)
}

func (s Shorts) Redirect(c echo.Context) error {
	short := c.Request().URL.Path[1:]

	url, err := s.service.GetURL(short)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, url)
}
