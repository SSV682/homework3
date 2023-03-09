package handlers

import (
	"delivery-service/internal/handlers/delivery"
	"delivery-service/internal/handlers/health"
	"delivery-service/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsEndpointName  = "/metrics"
	healthEndpointName   = "/health"
	productsEndpointName = "/products"
	DetailProductURL     = "/:product_id"
	ListURL              = ""
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

	h := delivery.NewHandler(rs.s)
	hh := health.NewHealth()

	e.GET(metricsEndpointName, echo.WrapHandler(promhttp.Handler()))
	e.GET(healthEndpointName, hh.Health)

	api := e.Group("/api")
	stableGroups := api.Group(VersionApi)
	{
		orders := stableGroups.Group(productsEndpointName)
		{
			orders.GET(ListURL, h.ListDeliveries)
		}
	}
	return nil
}
