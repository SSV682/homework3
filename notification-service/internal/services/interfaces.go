package services

import (
	"context"
	domain "notification-service/internal/domain/models"
)

type Service interface {
	List(ctx context.Context) ([]*domain.Notification, error)
}

type RunAsService interface {
	Run(ctx context.Context)
}
