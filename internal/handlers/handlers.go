package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"homework2/internal/handlers/health"
	"homework2/internal/handlers/users"
	promMiddleware "homework2/internal/middleware"
	"homework2/internal/services"
)

const (
	metricsEndpointName   = "/metrics"
	healthEndpointName    = "/health"
	usersEndpointName     = "/user"
	usersByIDEndpointName = usersEndpointName + "/:id"
)

const (
	VersionApi = "/v1"
)

type RegisterServices struct {
	s services.UserValueService
}

func NewRegisterServices(service services.UserValueService) *RegisterServices {
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

	stableGroups.GET(usersByIDEndpointName, h.GetUser)
	stableGroups.POST(usersEndpointName, h.CreateUser)
	stableGroups.PUT(usersByIDEndpointName, h.UpdateUser)
	stableGroups.DELETE(usersByIDEndpointName, h.DeleteUser)

	return nil
}
