package delivery

import (
	"context"
	domain "delivery-service/internal/domain/models"
	"delivery-service/internal/provider"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	deliveryApproved = "delivery_approved"
	deliveryRejected = "delivery_rejected"
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

	de := domain.DeliveryEntry{
		OrderID:      command.Order.ID,
		OrderContent: command.Order.OrderContent,
		Address:      command.Order.Address,
		Date:         command.Order.Date,
	}

	if err := p.storageProv.CheckPossibleDelivery(ctx, de); err != nil {
		log.Errorf("reject: %v", err)
		responseCommand.Command.Status = deliveryRejected
	} else {
		responseCommand.Command.Status = deliveryApproved
	}

	p.commandProducerProv.SendCommand(ctx, responseCommand)
}

func (p *Processor) rejectFunc(ctx context.Context, command domain.RequestCommand) {
	if err := p.storageProv.RejectDelivery(ctx, command.Order.ID); err == nil {
		responseCommand := domain.ResponseCommand{
			Topic: p.orderServiceTopic,
			Command: domain.Command{
				OrderID: command.Order.ID,
				Status:  deliveryRejected,
			},
		}
		p.commandProducerProv.SendCommand(ctx, responseCommand)
	} else {
		log.Errorf("failed reject order: %#v", command.Order)
	}
}
