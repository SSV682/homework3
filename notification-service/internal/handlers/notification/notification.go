package notification

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"notification-service/internal/services"
)

const (
	productIDPathArg = "product_id"
	limitQueryParam  = "limit"
	offsetQueryParam = "offset"
)

type handler struct {
	service services.Service
}

func NewHandler(s services.Service) *handler {
	return &handler{service: s}
}

func (h *handler) ListNotification(ctx echo.Context) error {

	//limit, err := queryParamsToUInt64(ctx.QueryParam(limitQueryParam), 10)
	//if err != nil {
	//	return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("wrong limit parameter: %v", err).Error()})
	//}
	//
	//offset, err := queryParamsToUInt64(ctx.QueryParam(offsetQueryParam), 0)
	//if err != nil {
	//	return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("wrong offst parameter: %v", err).Error()})
	//}

	//filter := dto.FilterProductDTO{
	//	Quantity: nil,
	//	Limit:    limit,
	//	Offset:   offset,
	//}

	ccx := ctx.Request().Context()

	orders, err := h.service.List(ccx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	result := make([]*ResponseNotification, 0, len(orders))

	for _, v := range orders {
		result = append(result, &ResponseNotification{
			ID:      v.ID,
			Mail:    v.Mail,
			Message: v.Message,
		})
	}

	res := ResponseNotifications{
		Total:   len(orders),
		Results: result,
	}

	return ctx.JSON(http.StatusOK, res)
}
