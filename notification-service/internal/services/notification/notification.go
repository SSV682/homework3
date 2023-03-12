package notification

import (
	"context"
	"fmt"
	domain "notification-service/internal/domain/models"
	"notification-service/internal/provider"
)

type notificationService struct {
	storageProv provider.StorageProvider
}

func NewNotificationService(s provider.StorageProvider) *notificationService {
	return &notificationService{
		storageProv: s,
	}
}

func (s *notificationService) List(ctx context.Context) ([]*domain.Notification, error) {
	res, err := s.storageProv.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get orders list: %v", err)
	}

	return res, nil
}
