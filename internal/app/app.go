package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"homework2/config"
	v1 "homework2/internal/controllers/http/v1"
	"homework2/pkg/httpserver"
	"homework2/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	e := echo.New()
	v1.RegisterRouters(e, l, cfg)

	httpServer := httpserver.New(e, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - singnal" + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run- httpServer.Notify:%w", err))
	}

	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run- httpServer.Shutdown:%w", err))
	}
}
