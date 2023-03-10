package services

import (
	domain "billing-service/internal/domain/models"
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Detail(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	Create(ctx context.Context, id uuid.UUID) error
	FillAccount(ctx context.Context, id uuid.UUID, amount float64) error
}

type RunAsService interface {
	Run(ctx context.Context)
}
