package services

import (
	"context"
	"user-service/internal/domain/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUser(ctx context.Context, token string) (models.User, error)
	DeleteUser(ctx context.Context, token string) error
	UpdateUser(ctx context.Context, token string, user *models.User) error
}
