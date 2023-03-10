package billing

import (
	"billing-service/internal/provider"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type OutboxConfig struct {
	OrderServiceTopic   string
	SystemBusTopic      string
	StorageProv         provider.StorageProvider
	CommandProducerProv provider.BrokerProducerProvider
}

type Outbox struct {
	startOnce           sync.Once
	orderServiceTopic   string
	systemBusTopic      string
	storageProv         provider.StorageProvider
	commandProducerProv provider.BrokerProducerProvider
}

func NewOutbox(cfg OutboxConfig) *Outbox {
	return &Outbox{
		orderServiceTopic:   cfg.OrderServiceTopic,
		systemBusTopic:      cfg.SystemBusTopic,
		storageProv:         cfg.StorageProv,
		commandProducerProv: cfg.CommandProducerProv,
	}
}

func (o *Outbox) Run(ctx context.Context) {
	o.startOnce.Do(func() {
		go o.start(ctx)
	})
	fmt.Println(fmt.Sprintf("outbox started"))
}

func (o *Outbox) start(ctx context.Context) func() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			for {
				command, err := o.storageProv.GetNextOutboxCommand(ctx)
				if err != nil || command == nil {
					break
				}

				log.Debugf("found: %#v, err: %s", command, err)
				if err = o.commandProducerProv.SendCommand(ctx, *command); err == nil {
					err := o.storageProv.DeleteOutboxCommand(ctx, command.ID)
					if err != nil {
						//TODO:
					}
				}
			}
		case <-ctx.Done():
			log.Infof("Contex faired! Stopping processor service...")
			break
		}
	}
}
