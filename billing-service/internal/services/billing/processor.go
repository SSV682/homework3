package billing

import (
	domain "billing-service/internal/domain/models"
	"billing-service/internal/provider"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	paymentApproved = "payment_approved"
	paymentRejected = "payment_rejected"
)

type Config struct {
	OrderServiceTopic   string
	SystemBusTopic      string
	StorageProv         provider.StorageProvider
	CommandConsumerProv provider.BrokerConsumerProvider
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
	de := domain.Order{
		ID:         command.Order.ID,
		UserID:     command.Order.UserID,
		TotalPrice: command.Order.TotalPrice,
	}

	if err := p.storageProv.CheckPossiblePayment(ctx, de); err != nil {
		log.Errorf("approve func: %v", err)
	}

}

func (p *Processor) rejectFunc(ctx context.Context, command domain.RequestCommand) {
	de := domain.Order{
		ID:         command.Order.ID,
		UserID:     command.Order.UserID,
		TotalPrice: command.Order.TotalPrice,
	}

	if err := p.storageProv.RejectPayment(ctx, de); err != nil {
		log.Errorf("failed reject order: %#v", command.Order)
	}
}
