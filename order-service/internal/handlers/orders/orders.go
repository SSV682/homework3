package orders

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
	"user-service/internal/domain/dto"
	"user-service/internal/services"
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
	payload := ctx.Request().Header.Get(tokenHeaderName)
	if payload == "" {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't cast x-jwt-token"})
	}

	userID, err := getUserID(payload)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: fmt.Errorf("couldn't get userID: %s", err).Error()})
	}

	idempotenceKey := ctx.Request().Header.Get(idempotenceKeyHeaderName)
	if idempotenceKey == "" {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't find x-request-id"})
	}

	var body CreateOrderRequest

	err = ctx.Bind(&body)
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

	orderDTO := order.OrderToDTO()

	return ctx.JSON(http.StatusOK, ResponseOrder{
		ID:         orderDTO.ID,
		UserID:     orderDTO.UserID,
		TotalPrice: orderDTO.TotalPrice,
		CreatedAt:  orderDTO.CreatedAt,
		Status:     orderDTO.Status,
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

	limit, err := queryParamsToUInt64(ctx.QueryParam(limitQueryParam))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("wrong limit parameter: %v", err).Error()})
	}

	offset, err := queryParamsToUInt64(ctx.QueryParam(offsetQueryParam))
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

	orderDTO := orders.OrdersToDTO()

	result := make([]*ResponseOrder, 0, len(orderDTO.Results))

	for _, v := range orderDTO.Results {
		result = append(result, &ResponseOrder{
			ID:         v.ID,
			UserID:     v.UserID,
			TotalPrice: v.TotalPrice,
			CreatedAt:  v.CreatedAt,
			Status:     v.Status,
		})
	}

	res := ResponseOrders{
		Total:   orderDTO.Total,
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

//func (h *handler) UpdateOrder(ctx echo.Context) error {
//	payload := ctx.Request().Header.Get(tokenHeaderName)
//	if payload == "" {
//		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't cast x-jwt-token"})
//	}
//
//	userID, err := getUserID(payload)
//	if err != nil {
//		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: fmt.Errorf("couldn't get userID: %s", err).Error()})
//	}
//
//	var body CreateOrderRequest
//
//	err = ctx.Bind(&body)
//	if err != nil {
//		return ctx.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
//	}
//
//	dto := dto.OrderRequestDTO{
//		UserID:         userID,
//		IdempotencyKey: idempotenceKey,
//		TotalPrice:     body.TotalPrice,
//		CreatedAt:      time.Now(),
//	}
//
//	cct := ctx.Request().Context()
//
//	err = h.service.Update(cct, , ,&dto)
//	if err != nil {
//		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
//	}
//
//	return ctx.JSON(http.StatusCreated, ResponseCreated{ID: id})
//}
