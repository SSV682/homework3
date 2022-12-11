package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"user-service/internal/handlers/health"
	"user-service/internal/handlers/users"
	promMiddleware "user-service/internal/middleware"
	"user-service/internal/services"
)

const (
	metricsEndpointName = "/metrics"
	healthEndpointName  = "/health"
	usersEndpointName   = "/user"
)

const (
	VersionApi = "/v1"
)

type RegisterServices struct {
	s services.UserService
}

func NewRegisterServices(service services.UserService) *RegisterServices {
	return &RegisterServices{s: service}
}

func RegisterHandlers(e *echo.Echo, rs *RegisterServices) error {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(promMiddleware.PrometheusMiddleware())

	h := users.NewHandler(rs.s)
	hh := health.NewHealth()

	e.GET(metricsEndpointName, echo.WrapHandler(promhttp.Handler()))

	e.GET(healthEndpointName, hh.Health)

	api := e.Group("/api")
	stableGroups := api.Group(VersionApi)

	stableGroups.GET(usersEndpointName, h.GetUser)
	stableGroups.POST(usersEndpointName, h.CreateUser)
	stableGroups.PATCH(usersEndpointName, h.UpdateUser)
	stableGroups.DELETE(usersEndpointName, h.DeleteUser)

	return nil
}
