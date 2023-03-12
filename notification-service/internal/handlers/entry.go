package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"notification-service/internal/handlers/health"
	"notification-service/internal/handlers/notification"
	"notification-service/internal/services"
)

const (
	metricsEndpointName       = "/metrics"
	healthEndpointName        = "/health"
	notificationsEndpointName = "/notifications"
	listURL                   = ""
)

const (
	VersionApi = "/v1"
)

type RegisterServices struct {
	s services.Service
}

func NewRegisterServices(service services.Service) *RegisterServices {
	return &RegisterServices{s: service}
}

func RegisterHandlers(e *echo.Echo, rs *RegisterServices) error {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := notification.NewHandler(rs.s)
	hh := health.NewHealth()

	e.GET(metricsEndpointName, echo.WrapHandler(promhttp.Handler()))
	e.GET(healthEndpointName, hh.Health)

	api := e.Group("/api")
	stableGroups := api.Group(VersionApi)
	{
		orders := stableGroups.Group(notificationsEndpointName)
		{
			orders.GET(listURL, h.ListNotification)
		}
	}
	return nil
}
