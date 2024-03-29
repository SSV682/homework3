package handlers

import (
	"auth-service/internal/handlers/auth"
	"auth-service/internal/handlers/health"
	"auth-service/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsEndpointName = "/metrics"
	healthEndpointName  = "/health"
	authEndpointName    = "/auth"
	loginEndpointName   = "/login"
	//signUpEndpointName  = "/signup"
	keysEndpointName = "/keys"
)

const (
	VersionApi = "/v1"
)

type RegisterServices struct {
	s services.AuthService
}

func NewRegisterServices(service services.AuthService) *RegisterServices {
	return &RegisterServices{s: service}
}

func RegisterHandlers(e *echo.Echo, rs *RegisterServices) error {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := auth.NewHandler(rs.s)
	hh := health.NewHealth()

	e.GET(metricsEndpointName, echo.WrapHandler(promhttp.Handler()))
	e.GET(healthEndpointName, hh.Health)

	api := e.Group("/api")
	stableGroups := api.Group(VersionApi)
	{
		stableGroups.POST(authEndpointName, h.CheckUser)
		stableGroups.POST(loginEndpointName, h.LoginUser)
		stableGroups.GET(keysEndpointName, h.Keys)
	}

	return nil
}
