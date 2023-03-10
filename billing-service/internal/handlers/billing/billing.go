package billing

import (
	"billing-service/internal/services"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	accountIDArg = "id"
)

type handler struct {
	service services.Service
}

func NewHandler(s services.Service) *handler {
	return &handler{service: s}
}

func (h *handler) DetailAccount(ctx echo.Context) error {

	paramID := ctx.Param(accountIDArg)

	var accountID uuid.UUID
	err := accountID.UnmarshalText([]byte(paramID))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("bad id parametr %v", err).Error()})
	}

	ccx := ctx.Request().Context()

	account, err := h.service.Detail(ccx, accountID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, ResponseAccount{
		ID:     account.UserID.String(),
		Amount: account.Amount,
	})
}

func (h *handler) CreateAccount(ctx echo.Context) error {
	paramID := ctx.Param(accountIDArg)

	var accountID uuid.UUID
	err := accountID.UnmarshalText([]byte(paramID))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("bad id parametr %v", err).Error()})
	}

	ccx := ctx.Request().Context()

	err = h.service.Create(ccx, accountID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h *handler) FillAccount(ctx echo.Context) error {

	var body RequestFillAccount

	err := ctx.Bind(&body)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	var accountID uuid.UUID
	err = accountID.UnmarshalText([]byte(body.ID))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: fmt.Errorf("bad id parametr %v", err).Error()})
	}

	ccx := ctx.Request().Context()

	err = h.service.FillAccount(ccx, accountID, body.Amount)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, nil)
}
