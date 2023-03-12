package provider

import (
	"context"
	domain "notification-service/internal/domain/models"
)

type StorageProvider interface {
	Create(ctx context.Context, p domain.Order) error
	List(ctx context.Context) ([]*domain.Notification, error)
}

type BrokerConsumerProvider interface {
	StartConsume(ctx context.Context) (<-chan domain.Order, <-chan error, error)
}
