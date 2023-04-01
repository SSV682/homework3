package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	domain "order-service/internal/domain/models"
	"order-service/internal/handlers/health"
	"order-service/internal/handlers/orders"
	"order-service/internal/services"
)

const (
	metricsEndpointName = "/metrics"
	healthEndpointName  = "/health"
	ordersEndpointName  = "/orders"
	DetailOrderURL      = "/:order_id"
	ListURL             = ""
	CancelURL           = "/:order_id/cancellation"
)

const (
	VersionName = "/v1"
	ApiName     = "/api"
)

type RegisterServices struct {
	service   services.OrderService
	validator domain.Validator
}

func NewRegisterServices(service services.OrderService, validator domain.Validator) *RegisterServices {
	return &RegisterServices{
		service:   service,
		validator: validator,
	}
}

func RegisterHandlers(e *echo.Echo, rs *RegisterServices) error {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := orders.NewHandler(rs.service, rs.validator)
	hh := health.NewHealth()

	e.GET(metricsEndpointName, echo.WrapHandler(promhttp.Handler()))
	e.GET(healthEndpointName, hh.Health)

	api := e.Group(ApiName)
	stableGroups := api.Group(VersionName).Group(ordersEndpointName)
	{
		stableGroups.POST(ListURL, h.CreateOrder)
		stableGroups.GET(DetailOrderURL, h.DetailOrder)
		stableGroups.GET(ListURL, h.ListOrder)
		stableGroups.DELETE(DetailOrderURL, h.DeleteOrder)
		stableGroups.PUT(CancelURL, h.CancelOrder)

	}

	return nil
}
