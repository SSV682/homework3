package health

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
}

func NewHealth() *handler {
	return &handler{}
}

func (impl *handler) Health(c echo.Context) error {
	o := Response{
		Status: "OK",
	}

	return c.JSON(http.StatusOK, o)
}
