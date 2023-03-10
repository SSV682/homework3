package auth

import (
	"auth-service/internal/domain/errors"
	"auth-service/internal/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
	service services.AuthService
}

func NewHandler(s services.AuthService) *handler {
	return &handler{service: s}
}

func (h *handler) LoginUser(ctx echo.Context) error {
	var user User

	err := echo.QueryParamsBinder(ctx).
		String("username", &user.Username).
		String("password", &user.Password).
		BindError()

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	cct := ctx.Request().Context()
	token, err := h.service.LoginUser(cct, user.Username, user.Password)
	if err != nil {
		return ctx.JSON(getStatusCode(err), Response{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, TokenResponse{AccessToken: *token})
}

func (h *handler) CheckUser(ctx echo.Context) error {
	var token Token
	err := ctx.Bind(&token)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{Message: fmt.Errorf("couldn't bind token: %s", err).Error()})
	}

	cct := ctx.Request().Context()
	claims, err := h.service.CheckUser(cct, token.Token)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	response := &UserResponse{
		UserID: claims.ID,
	}
	return ctx.JSON(http.StatusOK, response)
}

func (h *handler) Keys(ctx echo.Context) error {
	keys, err := h.service.GetKeys()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, keys)
}

func getStatusCode(err error) int {
	if err != nil {
		return http.StatusOK
	}
	switch err {
	case errors.ErrConflict:
		return http.StatusInternalServerError
	case errors.ErrNonExistentId:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
