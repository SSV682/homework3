package services

import (
	"context"
	"stock-service/internal/domain/dto"
	domain "stock-service/internal/domain/models"
)

type Service interface {
	Create(ctx context.Context, dto dto.ProductRequestDTO) (int64, error)
	List(ctx context.Context, filter dto.FilterProductDTO) ([]*domain.Product, error)
	Detail(ctx context.Context, productID int64) (*domain.Product, error)
	Update(ctx context.Context, productID int64, order domain.Product) error
	Fill(ctx context.Context, request dto.FillRequestDTO) error
}

type RunAsService interface {
	Run(ctx context.Context)
}
