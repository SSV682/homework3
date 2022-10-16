package app

import (
	"fmt"
	"github.com/dimfeld/httptreemux"
	"homework/config"
	v1 "homework/internal/controllers/http/v1"
	"homework/pkg/httpserver"
	"homework/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	router := httptreemux.New()
	v1.RegisterRouters(router)
	httpServer := httpserver.New(router, httpserver.Port(cfg.Port))

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
