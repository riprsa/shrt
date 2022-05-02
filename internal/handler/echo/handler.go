package handler

import (
	"github.com/hararudoka/shrt/internal/service"
	"github.com/labstack/echo/v4"
)

type handler struct {
	*service.Service
}

func New(service *service.Service) *handler {
	return &handler{
		service,
	}
}

type Router interface {
	Register(h handler, g *echo.Group)
}

func (h handler) NewGroup(g *echo.Group, r Router) {
	r.Register(h, g)
}
