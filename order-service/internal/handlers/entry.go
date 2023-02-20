package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	VersionApi     = "/v1"
	orderIDPathArg = "order_id"
)

type RegisterServices struct {
	s services.OrderService
}

func NewRegisterServices(service services.OrderService) *RegisterServices {
	return &RegisterServices{s: service}
}

func RegisterHandlers(e *echo.Echo, rs *RegisterServices) error {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := orders.NewHandler(rs.s)
	hh := health.NewHealth()

	e.GET(metricsEndpointName, echo.WrapHandler(promhttp.Handler()))
	e.GET(healthEndpointName, hh.Health)

	api := e.Group("/api")
	stableGroups := api.Group(VersionApi)
	{
		orders := stableGroups.Group(ordersEndpointName)
		{
			orders.POST(ListURL, h.CreateOrder)
			orders.GET(DetailOrderURL, h.DetailOrder)
			orders.GET(ListURL, h.ListOrder)
			orders.DELETE(DetailOrderURL, h.DeleteOrder)
			orders.PUT(CancelURL, h.CancelOrder)
		}
	}
	return nil
}
