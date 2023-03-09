package provider

import (
	"context"
	domain "delivery-service/internal/domain/models"
	"time"
)

type StorageProvider interface {
	ListDelivery(ctx context.Context, date time.Time) ([]*domain.DeliveryEntry, error)
	CheckPossibleDelivery(ctx context.Context, entry domain.DeliveryEntry) error
	RejectDelivery(ctx context.Context, orderID int64) error
}

type BrokerConsumerProvider interface {
	StartConsume(ctx context.Context) (<-chan domain.RequestCommand, <-chan error, error)
}

type BrokerProducerProvider interface {
	//SendMessage(ctx context.Context, command domain.Message) error
	SendCommand(ctx context.Context, command domain.ResponseCommand) error
}
