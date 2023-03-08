package provider

import (
	"context"
	"stock-service/internal/domain/dto"
	domain "stock-service/internal/domain/models"
)

type StorageProvider interface {
	CreateProduct(ctx context.Context, p domain.Product) (int64, error)
	UpdateProduct(ctx context.Context, id int64, p domain.Product) error
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	RavageStock(ctx context.Context, productsOrder []domain.Product) error
	FillStock(ctx context.Context, productsOrder []domain.Product) error
	ListStock(ctx context.Context, filter dto.FilterProductDTO) ([]*domain.Product, error)
}

type BrokerConsumerProvider interface {
	StartConsume(ctx context.Context) (<-chan domain.RequestCommand, <-chan error, error)
}

type BrokerProducerProvider interface {
	//SendMessage(ctx context.Context, command domain.Message) error
	SendCommand(ctx context.Context, command domain.ResponseCommand) error
}
