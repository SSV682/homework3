package orchestrator

import (
	"context"
	log "github.com/sirupsen/logrus"
	domain "order-service/internal/domain/models"
	"order-service/internal/provider"
	"sync"
)

type Config struct {
	CommandCh            <-chan domain.OrderCommand
	UpdateCh             chan<- domain.OrderCommand
	BillingServiceTopic  string
	StockServiceTopic    string
	DeliveryServiceTopic string
	CommandConsumerProv  provider.BrokerConsumerProvider
	CommandProducerProv  provider.BrokerProducerProvider
}

type Orchestrator struct {
	startOnce       sync.Once
	commandCh       <-chan domain.OrderCommand //для сервиса пользователей
	commandSourceCh <-chan domain.OrderCommand // для кафки
	UpdateCh        chan<- domain.OrderCommand

	billingServiceTopic  string
	stockServiceTopic    string
	deliveryServiceTopic string
	commandConsumerProv  provider.BrokerConsumerProvider
	commandProducerProv  provider.BrokerProducerProvider
	sagas                *domain.SagaSet
}

func NewOrchestrator(cfg Config) *Orchestrator {
	return &Orchestrator{
		commandCh:            cfg.CommandCh,
		UpdateCh:             cfg.UpdateCh,
		billingServiceTopic:  cfg.BillingServiceTopic,
		stockServiceTopic:    cfg.StockServiceTopic,
		deliveryServiceTopic: cfg.DeliveryServiceTopic,
		commandConsumerProv:  cfg.CommandConsumerProv,
		commandProducerProv:  cfg.CommandProducerProv,
		sagas:                domain.NewSagaSet(),
	}
}

func (o *Orchestrator) Register(order *domain.Order) {
	saga := domain.NewSaga(order, o.billingServiceTopic, o.stockServiceTopic, o.deliveryServiceTopic)
	o.sagas.Register(order.ID, saga)
}

func (o *Orchestrator) Run(ctx context.Context) {
	payloadCh, _, err := o.commandConsumerProv.StartConsume(ctx)
	if err != nil {
		log.Errorf("failed consumer: %v", err)
	}

	o.commandSourceCh = payloadCh

	o.startOnce.Do(func() {
		go o.start(ctx)
	})
}

func (o *Orchestrator) start(ctx context.Context) {
	for {
		select {
		case msg := <-o.commandCh: //создать, отменить
			o.executeCommand(ctx, msg)
		case msg := <-o.commandSourceCh: //подтверждение или отмена оплаты, подтверждение или отмена доставки
			o.executeCommand(ctx, msg)
		case <-ctx.Done():
			close(o.UpdateCh)
			log.Infof("Contex faired! Stopping orchestrator service...")
			break
		}
	}
}

func (o *Orchestrator) executeCommand(ctx context.Context, command domain.OrderCommand) {
	saga := o.sagas.Get(command.OrderID)
	if saga == nil {
		return
	}

	step := saga.NextState(command)

	switch step.Action {
	case domain.NextStep, domain.Retry:
		if err := o.commandProducerProv.SendCommand(ctx, step.Command); err != nil {
			log.Errorf("send command %v failed: %v", step.Command, err)
		}

		o.UpdateCh <- domain.OrderCommand{
			OrderID: command.OrderID,
			Status:  step.Status,
		}

	case domain.End:
		o.UpdateCh <- domain.OrderCommand{
			OrderID: command.OrderID,
			Status:  step.Status,
		}
		o.sagas.Remove(command.OrderID)
	case domain.Inaction:
	}
}
