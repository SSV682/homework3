package delivery

import (
	"delivery-service/internal/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	deliveryATQueryParam = "delivery_at"
)

type handler struct {
	service services.Service
}

func NewHandler(s services.Service) *handler {
	return &handler{service: s}
}

func (h *handler) ListDeliveries(ctx echo.Context) error {

	deliveryAt, err := queryParamsToUInt64(ctx.QueryParam(deliveryATQueryParam), 10)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("wrong limit parameter: %v", err).Error()})
	}

	ccx := ctx.Request().Context()

	orders, err := h.service.List(ccx, time.Unix(int64(deliveryAt), 0))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	result := make([]*ResponseDelivery, 0, len(orders))

	for _, v := range orders {
		result = append(result, &ResponseDelivery{
			ID:           v.ID,
			OrderId:      v.OrderID,
			OrderContent: v.OrderContent,
			Address:      v.Address,
			DeliveryAt:   v.Date.Unix(),
		})
	}

	res := ResponseDeliveries{
		Total:   len(orders),
		Results: result,
	}

	return ctx.JSON(http.StatusOK, res)
}
