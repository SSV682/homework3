package auth

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strings"
	"user-service/internal/domain/errors"
	"user-service/internal/domain/models"
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

func (h *handler) SignUp(ctx echo.Context) error {
	var user models.User
	err := ctx.Bind(&user)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if ok, err := isRequestValid(&user); !ok {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	cct := ctx.Request().Context()
	i, err := h.service.SignUpUser(cct, &user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, i)
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

func isRequestValid(u *models.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return false, err
	}
	return true, nil
}
