package stock

import (
	"context"
	"fmt"
	"stock-service/internal/domain/dto"
	domain "stock-service/internal/domain/models"
	"stock-service/internal/provider"
)

type stockService struct {
	storageProv provider.StorageProvider
}

func NewStockService(s provider.StorageProvider) *stockService {
	return &stockService{
		storageProv: s,
	}
}

func (s *stockService) Create(ctx context.Context, request dto.ProductRequestDTO) (int64, error) {
	product := domain.RestoreProductFromRequest(request)

	id, err := s.storageProv.CreateProduct(ctx, product)
	if err != nil {
		return 0, fmt.Errorf("failed create order: %v", err)
	}

	return id, nil
}

func (s *stockService) List(ctx context.Context, filter dto.FilterProductDTO) ([]*domain.Product, error) {
	res, err := s.storageProv.ListStock(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed get orders list: %v", err)
	}

	return res, nil
}

func (s *stockService) Detail(ctx context.Context, productID int64) (*domain.Product, error) {
	res, err := s.storageProv.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed get detail: %v", err)
	}

	return res, nil
}

func (s *stockService) Update(ctx context.Context, productID int64, product domain.Product) error {
	err := s.storageProv.UpdateProduct(ctx, productID, product)
	if err != nil {
		return fmt.Errorf("failed update: %v", err)
	}

	return nil
}

func (s *stockService) Fill(ctx context.Context, request dto.FillRequestDTO) error {
	products := make([]domain.Product, 0, len(request.Data))

	for _, v := range request.Data {
		products = append(products, domain.Product{
			ID:       v.ID,
			Quantity: v.Quantity,
			Name:     v.Name,
		})
	}

	err := s.storageProv.FillStock(ctx, products)
	if err != nil {
		return fmt.Errorf("failed filling: %v", err)
	}

	return nil
}
