package billing

import (
	"context"
	domain "delivery-service/internal/domain/models"
	"delivery-service/internal/provider"
	"fmt"
	"time"
)

type deliveryService struct {
	storageProv provider.StorageProvider
}

func NewDeliveryService(s provider.StorageProvider) *deliveryService {
	return &deliveryService{
		storageProv: s,
	}
}

func (s *deliveryService) List(ctx context.Context, date time.Time) ([]*domain.DeliveryEntry, error) {
	res, err := s.storageProv.ListDelivery(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed get deliveries list: %v", err)
	}

	return res, nil
}
