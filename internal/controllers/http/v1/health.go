package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Response struct {
	Status string
}

type healthHeandler struct {
}

func NewHealthImpl(e *echo.Group) {
	impl := &healthHeandler{}

	e.GET("/health", impl.health)
}

func (impl *healthHeandler) health(c echo.Context) error {
	o := Response{"OK"}
	return c.JSON(http.StatusOK, o)
}
