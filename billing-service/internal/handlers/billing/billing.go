package billing

import (
	"delivery-service/internal/services"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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

	if deliveryAt == 0 {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("wrong limit parameter: %v", err).Error()})
	}

	deliveryAtTime := time.Unix(int64(deliveryAt), 0)

	log.Debugf("delivery at: %d", deliveryAt)
	log.Debugf("delivery at: %s", deliveryAtTime)
	ccx := ctx.Request().Context()

	orders, err := h.service.List(ccx, deliveryAtTime)
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
