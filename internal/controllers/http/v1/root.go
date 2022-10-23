package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"homework2/config"
	"homework2/internal/db"
	"homework2/internal/domain/services"
	"homework2/pkg/logger"
	"strconv"
	"time"
)

func RegisterRouters(handler *echo.Echo, l logger.Interface, cfg *config.Config) {
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.Connection.User,
		cfg.Connection.Password,
		cfg.Connection.Host,
		cfg.Connection.Port,
		cfg.Connection.Dbname)

	userStorage, err := db.NewPgUserRepository(dsn)
	if err != nil {
		l.Error(err)
	}

	duration, err := strconv.Atoi(cfg.Timeout.Duration)
	if err != nil {
		l.Error(err)
	}
	timeoutContext := time.Duration(duration) * time.Second
	us := services.NewUserService(userStorage, timeoutContext)

	h := handler.Group("/v1")
	{
		NewUserImpl(h, us, l)
		NewHealthImpl(h)
	}
}
