package users

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"homework2/internal/domain/errors"
	"homework2/internal/domain/models"
	"homework2/internal/services"
	"net/http"
	"strconv"
)

type Handler struct {
	service services.UserValueService
}

func NewHandler(s services.UserValueService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateUser(ctx echo.Context) error {
	var user models.User
	err := ctx.Bind(&user)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if ok, err := isRequestValid(&user); !ok {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	cct := ctx.Request().Context()
	err = h.service.CreateUser(cct, &user)
	if err != nil {
		return ctx.JSON(getStatusCode(err), err.Error())
	}
	return ctx.JSON(http.StatusCreated, user)
}

func (h *Handler) GetUser(ctx echo.Context) error {
	idU, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, errors.ErrIncorrectParams)
	}

	ccx := ctx.Request().Context()
	user, err := h.service.GetUser(ccx, int64(idU))
	if err != nil {
		return ctx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteUser(ctx echo.Context) error {
	idU, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, errors.ErrIncorrectParams)
	}

	ccx := ctx.Request().Context()
	err = h.service.DeleteUser(ccx, int64(idU))
	if err != nil {
		return ctx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent) //?
}

func (h *Handler) UpdateUser(ctx echo.Context) error {
	var user models.User
	idU, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, errors.ErrIncorrectParams)
	}

	err = ctx.Bind(&user)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&user); !ok {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
	ccx := ctx.Request().Context()
	err = h.service.UpdateUser(ccx, int64(idU), &user)
	if err != nil {
		ctx.JSON(getStatusCode(err), err.Error())
	}

	return ctx.JSON(http.StatusOK, user)
}

func isRequestValid(u *models.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return false, err
	}
	return true, nil
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
