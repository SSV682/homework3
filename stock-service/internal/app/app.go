package app

import (
	"context"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"stock-service/internal/config"
	"stock-service/internal/handlers"
	"stock-service/internal/provider/kafka"
	"stock-service/internal/provider/sql"
	"stock-service/internal/services/stock"
	"syscall"
	"time"
)

type App struct {
	cfg        *config.Config
	httpServer *http.Server
}

func NewApp(configPath string) *App {
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configs: %v", err)
	}

	pool, err := initDBPool(cfg.Databases.Postgres)
	if err != nil {
		log.Fatalf("Failed to init database pool: %v", err)
	}

	handler := echo.New()

	producerConfig := kafka.ProducerConfig{
		Brokers: cfg.Kafka.BrokerAddresses,
	}

	sqlProv := sql.NewSQLProductProvider(pool)
	commandProducerProv := kafka.NewBrokerProducer(producerConfig)
	commandConsumerProv := kafka.NewBrokerConsumer(cfg.Kafka.BrokerAddresses, cfg.Topics.StockTopic)

	processorConfig := stock.Config{
		OrderServiceTopic:   cfg.Topics.OrderTopic,
		SystemBusTopic:      cfg.Topics.SystemBus,
		StorageProv:         sqlProv,
		CommandConsumerProv: commandConsumerProv,
		CommandProducerProv: commandProducerProv,
	}

	processorService := stock.NewProcessor(processorConfig)
	processorService.Run(context.Background())

	stockService := stock.NewStockService(sqlProv)
	rs := handlers.NewRegisterServices(stockService)
	err = handlers.RegisterHandlers(handler, rs)
	if err != nil {
		log.Fatalf("Failed to register handlers: %v", err)
	}

	log.Info("App created")

	return &App{
		cfg: &cfg,
		httpServer: &http.Server{
			Handler:      handler,
			Addr:         net.JoinHostPort("", cfg.HTTP.Port),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed listen and serve http server: %v", err)
		}
	}()

	log.Info("App has been started")
	a.waitGracefulShutdown(ctx, cancel)
}

func (a *App) waitGracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM, os.Interrupt,
	)

	log.Infof("Caught signal %s. Shutting down...", <-quit)

	cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Errorf("Failed to shutdown http server: %v", err)
	}
}
