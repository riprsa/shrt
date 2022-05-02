package handler

import (
	"database/sql"
	"net/http"

	"github.com/hararudoka/shrt/internal/model"

	"github.com/labstack/echo/v4"
)

type ShortService struct {
	handler
}

func (s ShortService) Register(h handler, g *echo.Group) {
	s.handler = h

	g.GET("", s.GetURL)
	g.POST("", s.SaveURL)

	g.GET("*", s.Redirect)
}

func (s ShortService) GetURL(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, "need json POST request with field url")
}

// SaveURL is a handler for POST request.
func (s ShortService) SaveURL(c echo.Context) error {
	// get url from request
	var data model.Data
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// get short url
	short, err := s.URL2Hash(data.URL)
	if err != nil {

		if err.Error() == "url is broken" {
			return c.JSON(http.StatusBadRequest, err.Error())
		} else if err.Error() == sql.ErrNoRows.Error() {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, model.Response{Short: short})
}

func (s ShortService) Redirect(c echo.Context) error {
	return c.JSON(302, model.Response{Short: c.Param("*")})
}
