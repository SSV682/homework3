package v1

import (
	"github.com/labstack/echo/v4"
	validator "gopkg.in/go-playground/validator.v9"
	"homework2/internal/domain/errors"
	"homework2/internal/domain/models"
	"homework2/internal/domain/services"
	"homework2/pkg/logger"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserService services.UserServiceInterface
	Logger      logger.Interface
}

func NewUserImpl(e *echo.Group, u services.UserServiceInterface, l logger.Interface) {
	handler := &UserHandler{
		UserService: u,
		Logger:      l,
	}

	e.POST("/user", handler.CreateUser)
	e.DELETE("/user/:id", handler.DeleteUser)
	e.GET("/user/:id", handler.GetUserById)
	e.PUT("/user/:id", handler.UpdateUser)

}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if ok, err := isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	err = h.UserService.CreateUser(ctx, &user)
	if err != nil {
		return c.JSON(getStatusCode(err), err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserById(c echo.Context) error {
	idU, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.ErrIncorrectParams)
	}
	id := int64(idU)
	ctx := c.Request().Context()
	user, err := h.UserService.GetUser(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	idU, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.ErrIncorrectParams)
	}
	id := int64(idU)
	ctx := c.Request().Context()
	err = h.UserService.DeleteUser(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent) //?
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	idU, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.ErrIncorrectParams)
	}

	var user models.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	id := int64(idU)

	if ok, err := isRequestValid(&user); !ok {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	err = h.UserService.UpdateUser(ctx, id, &user)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, user)
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
