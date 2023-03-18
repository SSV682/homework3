package provider

import (
	"context"
	domain "notification-service/internal/domain/models"
)

type StorageProvider interface {
	Create(ctx context.Context, p domain.Notification) error
	List(ctx context.Context) ([]*domain.Notification, error)
	UpdateUserInfo(ctx context.Context, user domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}

type BrokerConsumerProvider interface {
	StartConsume(ctx context.Context, errCh chan error) (<-chan domain.Order, error)
	StartConsumeUserUpdate(ctx context.Context, errCh chan error) (<-chan domain.User, error)
}
