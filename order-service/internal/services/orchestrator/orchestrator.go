package orchestrator

import (
	"context"
	log "github.com/sirupsen/logrus"
	domain "order-service/internal/domain/models"
	"order-service/internal/provider"
	"sync"
)

type Config struct {
	CommandCh           <-chan domain.OrderCommand
	UpdateCh            chan<- domain.OrderCommand
	BillingServiceTopic string
	StockServiceTopic   string
	CommandConsumerProv provider.BrokerConsumerProvider
	CommandProducerProv provider.BrokerProducerProvider
}

type Orchestrator struct {
	startOnce       sync.Once
	commandCh       <-chan domain.OrderCommand //для сервиса пользователей
	commandSourceCh <-chan domain.OrderCommand // для кафки
	UpdateCh        chan<- domain.OrderCommand

	billingServiceTopic string
	stockServiceTopic   string
	commandConsumerProv provider.BrokerConsumerProvider
	commandProducerProv provider.BrokerProducerProvider
	sagas               map[int64]*domain.Saga
}

func NewOrchestrator(cfg Config) *Orchestrator {
	return &Orchestrator{
		commandCh:           cfg.CommandCh,
		UpdateCh:            cfg.UpdateCh,
		billingServiceTopic: cfg.BillingServiceTopic,
		stockServiceTopic:   cfg.StockServiceTopic,
		commandConsumerProv: cfg.CommandConsumerProv,
		commandProducerProv: cfg.CommandProducerProv,
	}
}

func (o *Orchestrator) Register(order *domain.Order) {
	saga := domain.NewSaga(order, o.billingServiceTopic, o.stockServiceTopic)
	o.sagas[order.ID] = saga
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
			log.Infof("Contex faired! Stopping orchestrator service...")
			break
		}
	}
}

func (o *Orchestrator) executeCommand(ctx context.Context, command domain.OrderCommand) {
	saga, ok := o.sagas[command.OrderID]
	if !ok {
		//TODO: обработать ошибку
	}

	step := saga.NextState(command)
	switch step.Action {
	case domain.NextStep, domain.Retry:
		log.Infof("command: %v", step.Command)
		if err := o.commandProducerProv.SendCommand(ctx, step.Command); err != nil {
			log.Errorf("send command %v failed: %v", step.Command, err)
		}

		o.UpdateCh <- domain.OrderCommand{
			OrderID: step.Command.Order.ID,
			Status:  step.Status,
		}

	case domain.End:
		log.Infof("end position")
		delete(o.sagas, command.OrderID)
	case domain.Inaction:
	}
}
