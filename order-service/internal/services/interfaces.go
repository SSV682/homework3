package services

import (
	"context"
	"order-service/internal/domain/dto"
	domain "order-service/internal/domain/models"
)

type OrderService interface {
	Create(ctx context.Context, dto *dto.OrderRequestDTO) (int64, error)
	List(ctx context.Context, filter *dto.FilterOrderDTO) ([]*domain.Order, error)
	Detail(ctx context.Context, orderID int64, userID string) (*domain.Order, error)
	Delete(ctx context.Context, orderID int64, userID string) error
	Update(ctx context.Context, orderID int64, userID string, order *domain.Order) error
}

type RunAsService interface {
	Run(ctx context.Context)
}
