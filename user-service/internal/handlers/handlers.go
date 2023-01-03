package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"user-service/internal/handlers/health"
	"user-service/internal/handlers/users"
	customMiddleware "user-service/internal/middleware"
	"user-service/internal/services"
)

const (
	metricsEndpointName = "/metrics"
	healthEndpointName  = "/health"
	usersEndpointName   = "/user"
	signUpEndpointName  = "/signup"
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
	e.Use(customMiddleware.PrometheusMiddleware())

	h := users.NewHandler(rs.s)
	hh := health.NewHealth()

	e.GET(metricsEndpointName, echo.WrapHandler(promhttp.Handler()))
	e.GET(healthEndpointName, hh.Health)

	api := e.Group("/api")
	stableGroups := api.Group(VersionApi)
	{
		stableGroups.Use(customMiddleware.AuthMiddleware())
		stableGroups.GET(usersEndpointName, h.GetUser)
		stableGroups.PATCH(usersEndpointName, h.UpdateUser)
		stableGroups.DELETE(usersEndpointName, h.DeleteUser)

		stableGroups.POST(signUpEndpointName, h.CreateUser)

	}
	return nil
}
