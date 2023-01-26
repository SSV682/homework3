package provider

import (
	"context"
	"user-service/internal/domain/dto"
	"user-service/internal/domain/models"
)

type OrderProvider interface {
	CreateOrder(ctx context.Context, order *domain.Order) (int64, error)
	DetailOrder(ctx context.Context, id int64, userID string) (*domain.Order, error)
	ListOrders(ctx context.Context, dto *dto.FilterOrderDTO) ([]*domain.Order, error)
	DeleteOrder(ctx context.Context, id int64, userID string) error
	UpdateOrder(ctx context.Context, id int64, userID string, order *domain.Order) error
}

type RedisProvider interface {
	Write(ctx context.Context, key string, value int64) error
	Exist(ctx context.Context, key string) (bool, error)
	Read(ctx context.Context, key string) (int64, error)
}
