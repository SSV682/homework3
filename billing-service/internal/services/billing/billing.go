package billing

import (
	domain "billing-service/internal/domain/models"
	"billing-service/internal/provider"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type deliveryService struct {
	storageProv provider.StorageProvider
}

func NewDeliveryService(s provider.StorageProvider) *deliveryService {
	return &deliveryService{
		storageProv: s,
	}
}

func (s *deliveryService) Detail(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	res, err := s.storageProv.DetailAccount(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed get detail account: %v", err)
	}

	return res, nil
}

func (s *deliveryService) Create(ctx context.Context, id uuid.UUID) error {
	err := s.storageProv.CreateAccount(ctx, id)
	if err != nil {
		return fmt.Errorf("failed create account: %v", err)
	}

	return nil
}

func (s *deliveryService) FillAccount(ctx context.Context, id uuid.UUID, amount float64) error {
	err := s.storageProv.FillAccount(ctx, id, amount)
	if err != nil {
		return fmt.Errorf("failed filling account: %v", err)
	}

	return nil
}
