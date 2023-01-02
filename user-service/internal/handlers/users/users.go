package users

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strings"
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
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	userID := claims.UserID
	//token := ctx.Request().Header.Get("Authorization")
	//jwtString := strings.Split(token, "Bearer ")[1]

	//ccx := ctx.Request().Context()
	//user, err := h.service.GetUser(ccx, jwtString)
	//if err != nil {
	//	return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	//}
	return ctx.JSON(http.StatusOK, userID)
}

func (h *handler) DeleteUser(ctx echo.Context) error {
	token := ctx.Request().Header.Get("Authorization")
	jwtString := strings.Split(token, "Bearer ")[1]

	ccx := ctx.Request().Context()
	err := h.service.DeleteUser(ccx, jwtString)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent) //?
}

func (h *handler) UpdateUser(ctx echo.Context) error {
	var user models.User
	token := ctx.Request().Header.Get("Authorization")
	jwtString := strings.Split(token, "Bearer ")[1]

	ccx := ctx.Request().Context()
	err := ctx.Bind(&user)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&user); !ok {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
	err = h.service.UpdateUser(ccx, jwtString, &user)
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

//func getStatusCode(err error) int {
//	if err != nil {
//		return http.StatusOK
//	}
//	switch err {
//	case errors.ErrConflict:
//		return http.StatusInternalServerError
//	case errors.ErrNonExistentId:
//		return http.StatusNotFound
//	default:
//		return http.StatusInternalServerError
//	}
//}
