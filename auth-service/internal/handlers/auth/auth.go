package auth

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
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

	log.Infof("%v", user)
	cct := ctx.Request().Context()
	token, err := h.service.LoginUser(cct, user.Username, user.Password)
	if err != nil {
		return ctx.JSON(getStatusCode(err), err.Error())
	}

	return ctx.JSON(http.StatusCreated, token)
}

//func (h *handler) CheckUser(ctx echo.Context) error {
//	var user User
//
//	err := ctx.Bind(&user)
//	if err != nil {
//		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
//	}
//
//	cct := ctx.Request().Context()
//	token, err := h.service.LoginUser(cct, user.Username, user.Password)
//	if err != nil {
//		return ctx.JSON(getStatusCode(err), err.Error())
//	}
//
//	return ctx.JSON(http.StatusCreated, token)
//}

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
