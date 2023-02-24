package saga

import (
	"context"
	log "github.com/sirupsen/logrus"
	"order-service/internal/domain/dto"
	domain "order-service/internal/domain/models"
	"order-service/internal/provider"
	"sync"
)

type OrchestratorConfig struct {
	CommandCh           <-chan dto.OrderCommandDTO
	BillingServiceTopic string
	StockServiceTopic   string
	CommandConsumerProv provider.BrokerConsumerProvider
	CommandProducerProv provider.BrokerProducerProvider
	SqlProv             provider.OrderProvider
}

type Orchestrator struct {
	startOnce       sync.Once
	commandCh       <-chan dto.OrderCommandDTO //для сервиса пользователей
	commandSourceCh <-chan dto.OrderCommandDTO // для кафки

	billingServiceTopic string
	stockServiceTopic   string
	commandConsumerProv provider.BrokerConsumerProvider
	commandProducerProv provider.BrokerProducerProvider
	sqlProv             provider.OrderProvider
}

func NewOrchestrator(cfg *OrchestratorConfig) *Orchestrator {
	return &Orchestrator{
		commandCh:           cfg.CommandCh,
		billingServiceTopic: cfg.BillingServiceTopic,
		stockServiceTopic:   cfg.StockServiceTopic,
		commandConsumerProv: cfg.CommandConsumerProv,
		commandProducerProv: cfg.CommandProducerProv,
		sqlProv:             cfg.SqlProv,
	}
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

func (o *Orchestrator) executeCommand(ctx context.Context, command dto.OrderCommandDTO) {
	switch command.Status {
	case dto.Created:
		o.approvePayment(ctx, command.OrderID)
	case dto.Canceling:
		o.rejectStock(ctx, command.OrderID)
	case dto.PaymentApproved:
		o.approveStock(ctx, command.OrderID)
	case dto.StockApproved:
		o.approveOrder(ctx, command.OrderID)
	case dto.PaymentRejected:
		o.cancelOrder(ctx, command.OrderID)
	case dto.StockRejected:
		o.rejectPayment(ctx, command.OrderID)
	}
}

func (o *Orchestrator) approvePayment(ctx context.Context, id int64) {
	order, ok, err := o.sqlProv.GetOrderByIDThenUpdate(ctx, id, UpdateOrderStatusFunc(domain.PaymentPending))
	if err != nil {
		log.Errorf("sql approve payment failed: %v", err)
	}

	if ok {
		cm := dto.CommandDTO{
			Topic:       o.billingServiceTopic,
			CommandType: dto.Approve,
			Order:       *order.OrderToDTO(),
		}

		log.Infof("command: %v", cm)
		if err = o.commandProducerProv.SendMessage(ctx, cm); err != nil {
			log.Errorf("send message apporve payment failed: %v", err)
		}
	} else {

		o.cancelOrder(ctx, id)
	}

}

func (o *Orchestrator) approveStock(ctx context.Context, id int64) {
	log.Info("approve stock")
	order, ok, err := o.sqlProv.GetOrderByIDThenUpdate(ctx, id, UpdateOrderStatusFunc(domain.StockPending))
	if err != nil {
		log.Errorf("sql approve stock failed: %v", err)
	}

	cm := dto.CommandDTO{
		Order: *order.OrderToDTO(),
	}

	if ok {
		cm.Topic = o.stockServiceTopic
		cm.CommandType = dto.Approve
	} else {
		cm.Topic = o.billingServiceTopic
		cm.CommandType = dto.Reject
	}

	log.Infof("command: %v", cm)
	if err = o.commandProducerProv.SendMessage(ctx, cm); err != nil {
		log.Errorf("send message approve stock failed: %v", err)
	}
}

func (o *Orchestrator) rejectPayment(ctx context.Context, id int64) {
	order, _, err := o.sqlProv.GetOrderByIDThenUpdate(ctx, id, UpdateOrderStatusFunc(domain.PaymentRejecting))
	if err != nil {
		log.Errorf("sql reject payment failed: %v", err)
	}

	cm := dto.CommandDTO{
		Topic:       o.billingServiceTopic,
		CommandType: dto.Reject,
		Order:       *order.OrderToDTO(),
	}

	if err = o.commandProducerProv.SendMessage(ctx, cm); err != nil {
		log.Errorf("send message reject payment failed: %v", err)
	}
}

func (o *Orchestrator) approveOrder(ctx context.Context, id int64) {
	order, ok, err := o.sqlProv.GetOrderByIDThenUpdate(ctx, id, UpdateOrderStatusFunc(domain.Success))
	if err != nil {
		log.Errorf("approve order failed: %v", err)
	}

	if ok {
		cm := dto.CommandDTO{
			Topic:       o.stockServiceTopic,
			Order:       *order.OrderToDTO(),
			CommandType: dto.Reject,
		}

		log.Infof("command: %v", cm)
		if err = o.commandProducerProv.SendMessage(ctx, cm); err != nil {
			log.Errorf("send message reject payment failed: %v", err)
		}
	}

	//TODO: if not ok
}

func (o *Orchestrator) cancelOrder(ctx context.Context, id int64) {
	_, _, err := o.sqlProv.GetOrderByIDThenUpdate(ctx, id, UpdateOrderStatusFunc(domain.Canceled))
	if err != nil {
		log.Errorf("cancel order failed: %v", err)
	}
}

func (o *Orchestrator) rejectStock(ctx context.Context, id int64) {
	order, _, err := o.sqlProv.GetOrderByIDThenUpdate(ctx, id, UpdateOrderStatusFunc(domain.StockReject))
	if err != nil {
		log.Errorf("sql reject stock failed: %v", err)
	}

	cm := dto.CommandDTO{
		Topic:       o.stockServiceTopic,
		CommandType: dto.Reject,
		Order:       *order.OrderToDTO(),
	}

	if err = o.commandProducerProv.SendMessage(ctx, cm); err != nil {
		log.Errorf("send message reject stock failed: %v", err)
	}
}

func UpdateOrderStatusFunc(status domain.Status) domain.IntermediateOrderFunc {
	return func(o *domain.Order) (bool, error) {
		if o.Status() != domain.Canceling || status == domain.Canceled {
			o.SetStatus(status)
			log.Infof("order: %v", o)
		} else {
			return false, nil
		}

		return true, nil
	}
}