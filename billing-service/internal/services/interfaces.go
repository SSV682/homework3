package services

import (
	"context"
	domain "delivery-service/internal/domain/models"
	"time"
)

type Service interface {
	List(ctx context.Context, date time.Time) ([]*domain.DeliveryEntry, error)
}

type RunAsService interface {
	Run(ctx context.Context)
}
