package app

import (
	"context"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"order-service/internal/config"
	domain "order-service/internal/domain/models"
	"order-service/internal/handlers"
	"order-service/internal/provider/kafka"
	"order-service/internal/provider/redis"
	"order-service/internal/provider/sql"
	"order-service/internal/services/orchestrator"
	"order-service/internal/services/orders"
	updater "order-service/internal/services/updater"
	"os"
	"os/signal"
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

	client, err := initRedis(cfg.Cache)
	if err != nil {
		log.Fatalf("Failed to init redis: %v", err)
	}

	handler := echo.New()

	commandCh := make(chan domain.OrderCommand, 1000)
	updateCh := make(chan domain.OrderCommand, 1000)
	orderProv := sql.NewSQLBusinessRulesProvider(pool)
	redisProv := redis.NewRedisProvider(client)
	orderService := orders.NewOrdersService(orderProv, redisProv, commandCh)

	producerConfig := kafka.ProducerConfig{
		//Username: cfg.Kafka.SASL.Username,
		//Password: cfg.Kafka.SASL.Password,
		Brokers: cfg.Kafka.BrokerAddresses,
	}

	commandProducerProv := kafka.NewBrokerProducer(producerConfig)

	commandConsumerProv := kafka.NewBrokerConsumer(cfg.Kafka.BrokerAddresses, cfg.Topics.OrderTopic)

	sagaCfg := orchestrator.Config{
		CommandCh:           commandCh,
		UpdateCh:            updateCh,
		BillingServiceTopic: cfg.Topics.BillingTopic,
		StockServiceTopic:   cfg.Topics.StockTopic,
		CommandConsumerProv: commandConsumerProv,
		CommandProducerProv: commandProducerProv,
	}

	orchestrator := orchestrator.NewOrchestrator(sagaCfg)
	orchestrator.Run(context.Background())

	updaterCfg := updater.Config{
		CommandCh:           updateCh,
		SystemBusTopic:      cfg.Topics.SystemBus,
		StorageProv:         orderProv,
		CommandProducerProv: commandProducerProv,
	}
	updater := updater.NewUpdater(updaterCfg)
	updater.Run(context.Background())

	rs := handlers.NewRegisterServices(orderService)
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
