package users

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"user-service/internal/domain/models"
	"user-service/internal/services"
)

type handler struct {
	service services.UserService
}

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewHandler(s services.UserService) *handler {
	return &handler{service: s}
}

func (h *handler) CreateUser(ctx echo.Context) error {
	var user models.User
	err := ctx.Bind(&user)
	log.Infof("%v", user)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if ok, err := isRequestValid(&user); !ok {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	cct := ctx.Request().Context()
	i, err := h.service.CreateUser(cct, &user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, i)
}

func (h *handler) GetUser(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)

	ccx := ctx.Request().Context()
	user, err := h.service.GetUser(ccx, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, user)
}

func (h *handler) DeleteUser(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)

	ccx := ctx.Request().Context()
	err := h.service.DeleteUser(ccx, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent) //?
}

func (h *handler) UpdateUser(ctx echo.Context) error {
	var user models.User

	userID := ctx.Get("userID").(string)

	ccx := ctx.Request().Context()
	err := ctx.Bind(&user)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&user); !ok {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
	err = h.service.UpdateUser(ccx, userID, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
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
