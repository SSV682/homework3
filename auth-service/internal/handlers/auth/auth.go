package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"user-service/internal/domain/errors"
	"user-service/internal/services"
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
	//err := ctx.Bind(&user)
	//(&DefaultBinder{}).
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	cct := ctx.Request().Context()
	token, err := h.service.LoginUser(cct, user.Username, user.Password)
	if err != nil {
		return ctx.JSON(getStatusCode(err), err.Error())
	}

	return ctx.JSON(http.StatusCreated, token)
}

func (h *handler) CheckUser(ctx echo.Context) error {
	token := ctx.Request().Header.Get("Authorization")
	jwtString := strings.Split(token, "Bearer ")[1]
	cct := ctx.Request().Context()
	ok, err := h.service.CheckUser(cct, jwtString)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, err.Error())
	}
	if ok {
		return ctx.JSON(http.StatusOK, "ok")
	}
	return ctx.JSON(http.StatusUnauthorized, "invalid token")
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
