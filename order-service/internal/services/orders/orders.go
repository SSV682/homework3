package orders

import (
	"context"
	"fmt"
	"order-service/internal/domain/dto"
	domain "order-service/internal/domain/models"
	"order-service/internal/provider"
	"order-service/internal/services/orchestrator"
)

type orderService struct {
	orchestrator *orchestrator.Orchestrator
	storProv     provider.StorageProvider
	tempProv     provider.TempStorageProvider
	commandCh    chan domain.OrderCommand
}

func NewOrdersService(
	s provider.StorageProvider,
	t provider.TempStorageProvider,
	commandCh chan domain.OrderCommand,
	o *orchestrator.Orchestrator,
) *orderService {
	return &orderService{
		orchestrator: o,
		storProv:     s,
		tempProv:     t,
		commandCh:    commandCh,
	}
}

func (o *orderService) Create(ctx context.Context, request *dto.OrderRequestDTO) (int64, error) {
	exist, err := o.tempProv.Exist(ctx, key(request.IdempotencyKey, request.UserID))
	if err != nil {
		return 0, fmt.Errorf("exist: %v", err)
	}

	if exist {
		if id, err := o.tempProv.Read(ctx, key(request.IdempotencyKey, request.UserID)); err == nil {
			return id, nil
		}

		return 0, domain.ErrIdempotencyConflict
	}

	order := domain.NewOrderFromDTO(request)

	id, err := o.storProv.CreateOrder(ctx, order)
	if err != nil {
		return 0, fmt.Errorf("create order: %v", err)
	}

	order.ID = id
	o.orchestrator.Register(order)

	err = o.tempProv.Write(ctx, key(request.IdempotencyKey, request.UserID), id)
	if err != nil {
		return 0, fmt.Errorf("save key: %v", err)
	}

	command := domain.OrderCommand{
		OrderID: id,
		Status:  domain.Created,
	}

	o.commandCh <- command

	return id, nil
}

func (o *orderService) Detail(ctx context.Context, orderID int64, userID string) (*domain.Order, error) {
	res, err := o.storProv.DetailOrder(ctx, orderID, userID)
	if err != nil {
		return nil, fmt.Errorf("detail: %v", err)
	}

	return res, nil
}

func (o *orderService) Delete(ctx context.Context, orderID int64, userID string) error {
	if err := o.storProv.DeleteOrder(ctx, orderID, userID); err != nil {
		return fmt.Errorf("delete: %v", err)
	}

	return nil
}

func (o *orderService) List(ctx context.Context, filter *dto.FilterOrderDTO) ([]*domain.Order, error) {
	res, err := o.storProv.ListOrders(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("orders list: %v", err)
	}

	return res, nil
}

func (o *orderService) Update(ctx context.Context, orderID int64, userID string, order *domain.Order) error {
	err := o.storProv.UpdateOrder(ctx, orderID, userID, order)
	if err != nil {
		return fmt.Errorf("update: %v", err)
	}

	return nil
}

func (o *orderService) Cancel(ctx context.Context, orderID int64, userID string) error {
	order, err := o.storProv.DetailOrder(ctx, orderID, userID)
	if err != nil {
		return fmt.Errorf("detail order: %v", err)
	}

	o.orchestrator.Register(order)
	command := domain.OrderCommand{
		OrderID: orderID,
		Status:  domain.Canceling,
	}

	o.commandCh <- command

	return nil
}
