package health

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type handler struct {
}

type Response struct {
	Status string
}

func NewHealth() *handler {
	return &handler{}
}

func (impl *handler) Health(c echo.Context) error {
	log.Info("Health check")
	o := Response{
		Status: "OK",
	}
	return c.JSON(http.StatusOK, o)
}
