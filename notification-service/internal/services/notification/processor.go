package notification

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	domain "notification-service/internal/domain/models"
	"notification-service/internal/provider"
	"sync"
)

type Config struct {
	SystemBusTopic      string
	StorageProv         provider.StorageProvider
	CommandConsumerProv provider.BrokerConsumerProvider
}

type Processor struct {
	startOnce           sync.Once
	commandCh           <-chan domain.Order
	orderServiceTopic   string
	systemBusTopic      string
	storageProv         provider.StorageProvider
	commandConsumerProv provider.BrokerConsumerProvider
}

func NewProcessor(cfg Config) *Processor {
	return &Processor{
		systemBusTopic:      cfg.SystemBusTopic,
		storageProv:         cfg.StorageProv,
		commandConsumerProv: cfg.CommandConsumerProv,
	}
}

func (p *Processor) Run(ctx context.Context) {
	payloadCh, _, err := p.commandConsumerProv.StartConsume(ctx)
	if err != nil {
		log.Errorf("failed consumer: %v", err)
	}

	p.commandCh = payloadCh

	p.startOnce.Do(func() {
		go p.start(ctx)
	})
	fmt.Println(fmt.Sprintf("processor started"))
}

func (p *Processor) start(ctx context.Context) func() {
	for {
		select {
		case command := <-p.commandCh:
			fmt.Println(fmt.Sprintf("command: %#v", command))
			p.executeCommand(ctx, command)
		case <-ctx.Done():
			log.Infof("Contex faired! Stopping processor service...")
			break
		}
	}
}

func (p *Processor) executeCommand(ctx context.Context, command domain.Order) {
	switch command.Status {
	case domain.Success, domain.Canceled:
		p.storageProv.Create(ctx, command)
	default:
	}
}
