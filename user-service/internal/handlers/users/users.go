package users

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"user-service/internal/domain/models"
	"user-service/internal/services"
)

const (
	tokenHeaderName = "x-jwt-token"
)

type handler struct {
	service services.UserService
}

func NewHandler(s services.UserService) *handler {
	return &handler{service: s}
}

func (h *handler) CreateUser(ctx echo.Context) error {
	var user models.User
	err := ctx.Bind(&user)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	if ok, err := isRequestValid(&user); !ok {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	cct := ctx.Request().Context()

	i, err := h.service.CreateUser(cct, &user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, i)
}

func (h *handler) GetUser(ctx echo.Context) error {
	payload, ok := ctx.Get(tokenHeaderName).(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't cast x-auth-token"})
	}

	userID, err := getUserID(payload)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}

	ccx := ctx.Request().Context()

	user, err := h.service.GetUser(ccx, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}

func (h *handler) DeleteUser(ctx echo.Context) error {
	payload, ok := ctx.Get(tokenHeaderName).(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't cast x-auth-token"})
	}

	userID, err := getUserID(payload)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}

	ccx := ctx.Request().Context()

	err = h.service.DeleteUser(ccx, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}
	return ctx.JSON(http.StatusNoContent, ResponseError{Message: "No content"})
}

func (h *handler) UpdateUser(ctx echo.Context) error {
	var user models.User

	payload, ok := ctx.Get(tokenHeaderName).(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: "couldn't cast x-auth-token"})
	}

	userID, err := getUserID(payload)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}

	ccx := ctx.Request().Context()
	err = ctx.Bind(&user)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	if ok, err := isRequestValid(&user); !ok {
		return ctx.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}
	err = h.service.UpdateUser(ccx, userID, &user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}
