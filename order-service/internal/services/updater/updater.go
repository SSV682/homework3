package updater

import (
	"context"
	log "github.com/sirupsen/logrus"
	domain "order-service/internal/domain/models"
	"order-service/internal/provider"
	"sync"
)

type Config struct {
	CommandCh           <-chan domain.OrderCommand
	SystemBusTopic      string
	StorageProv         provider.StorageProvider
	CommandProducerProv provider.BrokerProducerProvider
}

type Updater struct {
	startOnce           sync.Once
	commandCh           <-chan domain.OrderCommand
	systemBusTopic      string
	storageProv         provider.StorageProvider
	messageProducerProv provider.BrokerProducerProvider
}

func NewUpdater(cfg Config) *Updater {
	return &Updater{
		commandCh:           cfg.CommandCh,
		systemBusTopic:      cfg.SystemBusTopic,
		storageProv:         cfg.StorageProv,
		messageProducerProv: cfg.CommandProducerProv,
	}
}

func (u *Updater) Run(ctx context.Context) {
	u.startOnce.Do(func() {
		go u.start(ctx)
	})
}

func (u *Updater) start(ctx context.Context) {
	for command := range u.commandCh {
		order, err := u.storageProv.GetOrderByIDThenUpdate(ctx, command.OrderID, UpdateOrderStatusFunc(command.Status))
		if err != nil {
			log.Errorf("get then update: %v", err)
		}

		err = u.messageProducerProv.SendMessage(ctx, domain.Message{
			Topic: u.systemBusTopic,
			Order: *order,
		})
		if err != nil {
			log.Errorf("send message: %v", err)
		}
	}
}

func UpdateOrderStatusFunc(status domain.Status) domain.IntermediateOrderFunc {
	return func(o *domain.Order) (bool, error) {
		o.SetStatus(status)

		return true, nil
	}
}
