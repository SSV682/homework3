package provider

import (
	domain "billing-service/internal/domain/models"
	"context"
	"github.com/google/uuid"
)

type StorageProvider interface {
	CheckPossiblePayment(ctx context.Context, order domain.Order) error
	CreateOutboxCommand(ctx context.Context, command domain.ResponseCommand) (int64, error)
	GetNextOutboxCommand(ctx context.Context) (*domain.ReadyResponseCommand, error)
	DeleteOutboxCommand(ctx context.Context, id int64) error
	DetailAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	CreateAccount(ctx context.Context, id uuid.UUID) error
	FillAccount(ctx context.Context, id uuid.UUID, amount float64) error
	RejectPayment(ctx context.Context, order domain.Order) error
}

type BrokerConsumerProvider interface {
	StartConsume(ctx context.Context) (<-chan domain.RequestCommand, <-chan domain.Account, <-chan error, error)
}

type BrokerProducerProvider interface {
	//SendMessage(ctx context.Context, command domain.Message) error
	SendCommand(ctx context.Context, command domain.ReadyResponseCommand) error
}
