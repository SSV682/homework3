package orders

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"order-service/internal/domain/dto"
	domain "order-service/internal/domain/models"
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
	service   services.OrderService
	validInst domain.Validator
}

func NewHandler(s services.OrderService, v domain.Validator) *handler {
	return &handler{
		service:   s,
		validInst: v,
	}
}

func (h *handler) CreateOrder(ctx echo.Context) error {
	payload := ctx.Request().Header.Get(tokenHeaderName)
	if payload == "" {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "cast x-jwt-token"})
	}

	userID, err := getUserID(payload)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: fmt.Errorf("get userID: %s", err).Error()})
	}

	idempotenceKey := ctx.Request().Header.Get(idempotenceKeyHeaderName)
	if idempotenceKey == "" {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "find x-request-id"})
	}

	var body CreateOrderRequest

	if err = ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	if err = h.validInst.Struct(body); err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	orderDTO := dto.OrderRequestDTO{
		UserID:         userID,
		IdempotencyKey: idempotenceKey,
		TotalPrice:     body.TotalPrice,
		DeliveryAt:     body.DeliveryAt,
		Products:       body.Products,
		Address:        body.Address,
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
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "cast x-jwt-token"})
	}

	userID, err := getUserID(payload)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: fmt.Errorf("get userID: %s", err).Error()})
	}

	orderID := ctx.Param(orderIDPathArg)
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("id parametr %v", err).Error()})
	}

	ccx := ctx.Request().Context()

	order, err := h.service.Detail(ccx, int64(id), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOrderNotFound):
			return ctx.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
		default:
			return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		}
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
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "cast x-jwt-token"})
	}

	userID, err := getUserID(payload)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: fmt.Errorf("get userID: %s", err).Error()})
	}

	filter := dto.FilterOrderDTO{}
	filter.UserID = userID

	limit, err := queryParamsToUInt64(ctx.QueryParam(limitQueryParam), 10)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("limit parameter: %v", err).Error()})
	}

	offset, err := queryParamsToUInt64(ctx.QueryParam(offsetQueryParam), 0)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("offset parameter: %v", err).Error()})
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
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "cast x-jwt-token"})
	}

	userID, err := getUserID(payload)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: fmt.Errorf("get userID: %s", err).Error()})
	}

	orderID := ctx.Param(orderIDPathArg)
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("id parameter %v", err).Error()})
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
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("id parametr %v", err).Error()})
	}

	cct := ctx.Request().Context()

	err = h.service.Cancel(cct, int64(id), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOrderNotFound):
			return ctx.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
		default:
			return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		}
	}

	return ctx.JSON(http.StatusAccepted, ResponseCreated{ID: int64(id)})
}
