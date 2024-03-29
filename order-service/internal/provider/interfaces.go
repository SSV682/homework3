package provider

import (
	"context"
	"order-service/internal/domain/dto"
	"order-service/internal/domain/models"
)

type StorageProvider interface {
	CreateOrder(ctx context.Context, order *domain.Order) (int64, error)
	DetailOrder(ctx context.Context, id int64, userID string) (*domain.Order, error)
	GetOrderByID(ctx context.Context, id int64) (*domain.Order, error)
	ListOrders(ctx context.Context, dto *dto.FilterOrderDTO) ([]*domain.Order, error)
	DeleteOrder(ctx context.Context, id int64, userID string) error
	UpdateOrder(ctx context.Context, id int64, userID string, order *domain.Order) error
	GetOrderByIDThenUpdate(ctx context.Context, id int64, fn domain.IntermediateOrderFunc) (*domain.Order, error)
}

type TempStorageProvider interface {
	Write(ctx context.Context, key string, value int64) error
	Exist(ctx context.Context, key string) (bool, error)
	Read(ctx context.Context, key string) (int64, error)
}

type BrokerConsumerProvider interface {
	StartConsume(ctx context.Context) (<-chan domain.OrderCommand, <-chan error, error)
}

type BrokerProducerProvider interface {
	SendCommand(ctx context.Context, command domain.Command) error
	SendMessage(ctx context.Context, command domain.Message) error
}
