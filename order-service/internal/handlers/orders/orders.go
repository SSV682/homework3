package orders

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"order-service/internal/domain/dto"
	"order-service/internal/services"
	"strconv"
	"time"
)

const (
	tokenHeaderName          = "x-jwt-token"
	idempotenceKeyHeaderName = "X-Request-ID"
	orderIDPathArg           = "order_id"
	limitQueryParam          = "limit"
	offsetQueryParam         = "offset"
)

type handler struct {
	service services.OrderService
}

func NewHandler(s services.OrderService) *handler {
	return &handler{service: s}
}

func (h *handler) CreateOrder(ctx echo.Context) error {
	//payload := ctx.Request().Header.Get(tokenHeaderName)
	//if payload == "" {
	//	return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't cast x-jwt-token"})
	//}
	//
	//userID, err := getUserID(payload)
	//if err != nil {
	//	return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: fmt.Errorf("couldn't get userID: %s", err).Error()})
	//}

	userID := "503c3602-5c51-4848-b332-ead24b4e0621"

	idempotenceKey := ctx.Request().Header.Get(idempotenceKeyHeaderName)
	log.Info(idempotenceKey)
	if idempotenceKey == "" {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't find x-request-id"})
	}

	var body CreateOrderRequest

	err := ctx.Bind(&body)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	orderDTO := dto.OrderRequestDTO{
		UserID:         userID,
		IdempotencyKey: idempotenceKey,
		TotalPrice:     body.TotalPrice,
		CreatedAt:      time.Now(),
	}

	cct := ctx.Request().Context()

	id, err := h.service.Create(cct, &orderDTO)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, ResponseCreated{ID: id})
}

func (h *handler) DetailOrder(ctx echo.Context) error {
	payload := ctx.Request().Header.Get(tokenHeaderName)
	if payload == "" {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't cast x-jwt-token"})
	}

	userID, err := getUserID(payload)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: fmt.Errorf("couldn't get userID: %s", err).Error()})
	}

	orderID := ctx.Param(orderIDPathArg)
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("bad id parametr %v", err).Error()})
	}

	ccx := ctx.Request().Context()

	order, err := h.service.Detail(ccx, int64(id), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, ResponseOrder{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		Status:     string(order.Status),
	})
}

func (h *handler) ListOrder(ctx echo.Context) error {
	payload := ctx.Request().Header.Get(tokenHeaderName)
	if payload == "" {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't cast x-jwt-token"})
	}

	userID, err := getUserID(payload)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: fmt.Errorf("couldn't get userID: %s", err).Error()})
	}

	filter := dto.FilterOrderDTO{}
	filter.UserID = userID

	limit, err := queryParamsToUInt64(ctx.QueryParam(limitQueryParam), 10)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("wrong limit parameter: %v", err).Error()})
	}

	offset, err := queryParamsToUInt64(ctx.QueryParam(offsetQueryParam), 0)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("wrong offst parameter: %v", err).Error()})
	}

	filter.Limit = limit
	filter.Offset = offset

	ccx := ctx.Request().Context()

	orders, err := h.service.List(ccx, &filter)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	result := make([]*ResponseOrder, 0, len(orders))

	for _, v := range orders {
		result = append(result, &ResponseOrder{
			ID:         v.ID,
			UserID:     v.UserID,
			TotalPrice: v.TotalPrice,
			CreatedAt:  v.CreatedAt,
			Status:     string(v.Status),
		})
	}

	res := ResponseOrders{
		Total:   len(orders),
		Results: result,
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *handler) DeleteOrder(ctx echo.Context) error {
	payload := ctx.Request().Header.Get(tokenHeaderName)
	if payload == "" {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't cast x-jwt-token"})
	}

	userID, err := getUserID(payload)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: fmt.Errorf("couldn't get userID: %s", err).Error()})
	}

	orderID := ctx.Param(orderIDPathArg)
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("bad id parametr %v", err).Error()})
	}

	ccx := ctx.Request().Context()

	err = h.service.Delete(ccx, int64(id), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusNoContent, ResponseError{Message: "Resource deleted successfully"})
}

func (h *handler) CancelOrder(ctx echo.Context) error {
	//payload := ctx.Request().Header.Get(tokenHeaderName)
	//if payload == "" {
	//	return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't cast x-jwt-token"})
	//}
	//
	//userID, err := getUserID(payload)
	//if err != nil {
	//	return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: fmt.Errorf("couldn't get userID: %s", err).Error()})
	//}
	userID := "503c3602-5c51-4848-b332-ead24b4e0621"

	orderID := ctx.Param(orderIDPathArg)
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("bad id parametr %v", err).Error()})
	}

	cct := ctx.Request().Context()

	err = h.service.Cancel(cct, int64(id), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusAccepted, ResponseCreated{ID: int64(id)})
}
