package orders

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"user-service/internal/domain/dto"
	domain "user-service/internal/domain/models"
	"user-service/internal/provider"
)

type orderService struct {
	sqlProv   provider.OrderProvider
	redisProv provider.RedisProvider
}

func NewOrdersService(s provider.OrderProvider, t provider.RedisProvider) *orderService {
	return &orderService{
		sqlProv:   s,
		redisProv: t,
	}
}

func (o *orderService) Create(ctx context.Context, dto *dto.OrderRequestDTO) (int64, error) {
	exist, err := o.redisProv.Exist(ctx, key(dto.IdempotencyKey, dto.UserID))
	if err != nil {
		return 0, fmt.Errorf("failed check exist: %v", err)
	}

	if exist {
		if id, err := o.redisProv.Read(ctx, key(dto.IdempotencyKey, dto.UserID)); err == nil {
			return id, nil
		}

		return 0, errors.New("idempotency conflict")
	}

	id, err := o.sqlProv.CreateOrder(ctx, domain.NewOrderFromDTO(dto))
	if err != nil {
		return 0, fmt.Errorf("failed create order: %v", err)
	}

	err = o.redisProv.Write(ctx, key(dto.IdempotencyKey, dto.UserID), id)
	if err != nil {
		return 0, fmt.Errorf("failed save key: %v", err)
	}

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

func (o *orderService) List(ctx context.Context, filter *dto.FilterOrderDTO) (*domain.Orders, error) {
	res, err := o.sqlProv.ListOrders(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed get orders list: %v", err)
	}

	return domain.NewOrdersFromSlice(res), nil

}

func (o *orderService) Update(ctx context.Context, orderID int64, userID string, order *domain.Order) error {
	err := o.sqlProv.UpdateOrder(ctx, orderID, userID, order)
	if err != nil {
		return fmt.Errorf("failed update: %v", err)
	}

	return nil
}

func key(idempotenceKey, userID string) string {
	sb := strings.Builder{}

	sb.WriteString(idempotenceKey)
	sb.WriteString("_")
	sb.WriteString(userID)

	return sb.String()
}
