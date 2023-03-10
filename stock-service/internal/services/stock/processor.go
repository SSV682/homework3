package stock

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	domain "stock-service/internal/domain/models"
	"stock-service/internal/provider"
	"sync"
)

const (
	stockApproved = "stock_approved"
	stockRejected = "stock_rejected"
)

type Config struct {
	OrderServiceTopic   string
	SystemBusTopic      string
	StorageProv         provider.StorageProvider
	CommandConsumerProv provider.BrokerConsumerProvider
	CommandProducerProv provider.BrokerProducerProvider
}

type Processor struct {
	startOnce           sync.Once
	commandCh           <-chan domain.RequestCommand
	orderServiceTopic   string
	systemBusTopic      string
	storageProv         provider.StorageProvider
	commandConsumerProv provider.BrokerConsumerProvider
	commandProducerProv provider.BrokerProducerProvider
}

func NewProcessor(cfg Config) *Processor {
	return &Processor{
		orderServiceTopic:   cfg.OrderServiceTopic,
		systemBusTopic:      cfg.SystemBusTopic,
		storageProv:         cfg.StorageProv,
		commandConsumerProv: cfg.CommandConsumerProv,
		commandProducerProv: cfg.CommandProducerProv,
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
			p.executeCommand(ctx, command)
		case <-ctx.Done():
			log.Infof("Contex faired! Stopping processor service...")
			break
		}
	}
}

func (p *Processor) executeCommand(ctx context.Context, command domain.RequestCommand) {
	fmt.Println(fmt.Sprintf("command came: %#v", command))
	switch command.CommandType {
	case domain.Approve:
		p.approveFunc(ctx, command)
	case domain.Reject:
		p.rejectFunc(ctx, command)
	default:
	}
}

func (p *Processor) approveFunc(ctx context.Context, command domain.RequestCommand) {
	responseCommand := domain.ResponseCommand{
		Topic: p.orderServiceTopic,
		Command: domain.Command{
			OrderID: command.Order.ID,
		},
	}

	if err := p.storageProv.RavageStock(ctx, command.Order.Products); err != nil {
		responseCommand.Command.Status = stockRejected
	} else {
		responseCommand.Command.Status = stockApproved
	}

	p.commandProducerProv.SendCommand(ctx, responseCommand)
}

func (p *Processor) rejectFunc(ctx context.Context, command domain.RequestCommand) {
	if err := p.storageProv.FillStock(ctx, command.Order.Products); err == nil {
		responseCommand := domain.ResponseCommand{
			Topic: p.orderServiceTopic,
			Command: domain.Command{
				OrderID: command.Order.ID,
				Status:  stockRejected,
			},
		}
		p.commandProducerProv.SendCommand(ctx, responseCommand)
	} else {
		log.Errorf("failed reject order: %#v", command.Order)
	}
}
