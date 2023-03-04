package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"order-service/internal/domain/dto"
	domain "order-service/internal/domain/models"
	"order-service/internal/provider"
	"order-service/internal/services/orchestrator"
	"strings"
)

type orderService struct {
	orchestrator *orchestrator.Orchestrator
	sqlProv      provider.OrderProvider
	redisProv    provider.RedisProvider
	commandCh    chan domain.OrderCommand
}

func NewOrdersService(s provider.OrderProvider, t provider.RedisProvider, commandCh chan domain.OrderCommand) *orderService {
	return &orderService{
		sqlProv:   s,
		redisProv: t,
		commandCh: commandCh,
	}
}

func (o *orderService) Create(ctx context.Context, request *dto.OrderRequestDTO) (int64, error) {
	exist, err := o.redisProv.Exist(ctx, key(request.IdempotencyKey, request.UserID))
	if err != nil {
		return 0, fmt.Errorf("failed check exist: %v", err)
	}

	if exist {
		if id, err := o.redisProv.Read(ctx, key(request.IdempotencyKey, request.UserID)); err == nil {
			return id, nil
		}

		return 0, errors.New("idempotency conflict")
	}

	order := domain.NewOrderFromDTO(request)

	id, err := o.sqlProv.CreateOrder(ctx, order)
	if err != nil {
		return 0, fmt.Errorf("failed create order: %v", err)
	}

	order.ID = id
	o.orchestrator.Register(order)

	err = o.redisProv.Write(ctx, key(request.IdempotencyKey, request.UserID), id)
	if err != nil {
		return 0, fmt.Errorf("failed save key: %v", err)
	}

	command := domain.OrderCommand{
		OrderID: id,
		Status:  domain.Created,
	}

	o.commandCh <- command
	log.Infof("command sent: %v", command)

	return id, nil
}

func (o *orderService) Detail(ctx context.Context, orderID int64, userID string) (*domain.Order, error) {
	res, err := o.sqlProv.DetailOrder(ctx, orderID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed get detail: %v", err)
	}

	return res, nil
}

func (o *orderService) Delete(ctx context.Context, orderID int64, userID string) error {
	if err := o.sqlProv.DeleteOrder(ctx, orderID, userID); err != nil {
		return fmt.Errorf("failed delete: %v", err)
	}

	return nil
}

func (o *orderService) List(ctx context.Context, filter *dto.FilterOrderDTO) ([]*domain.Order, error) {
	res, err := o.sqlProv.ListOrders(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed get orders list: %v", err)
	}

	log.Infof("its ok, service: %v", res)

	return res, nil

}

func (o *orderService) Update(ctx context.Context, orderID int64, userID string, order *domain.Order) error {
	err := o.sqlProv.UpdateOrder(ctx, orderID, userID, order)
	if err != nil {
		return fmt.Errorf("failed update: %v", err)
	}

	return nil
}

func (o *orderService) Cancel(_ context.Context, orderID int64, _ string) error {
	log.Infof("order info when canceling: %d", orderID)
	command := domain.OrderCommand{
		OrderID: orderID,
		Status:  domain.Canceling,
	}

	log.Infof("command sent: %v", command)
	o.commandCh <- command
	return nil
}

func key(idempotenceKey, userID string) string {
	sb := strings.Builder{}

	sb.WriteString(idempotenceKey)
	sb.WriteString("_")
	sb.WriteString(userID)

	return sb.String()
}
