package health

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type healthHeandler struct {
}

func NewHealth() *healthHeandler {
	return &healthHeandler{}
}

func (impl *healthHeandler) Health(c echo.Context) error {
	log.Info("Health check")
	o := Response{
		Status: "OK",
	}
	return c.JSON(http.StatusOK, o)
}
