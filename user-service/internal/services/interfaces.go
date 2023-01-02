package services

import (
	"context"
	"user-service/internal/domain/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUser(ctx context.Context, userID string) (models.User, error)
	DeleteUser(ctx context.Context, userID string) error
	UpdateUser(ctx context.Context, userID string, user *models.User) error
}
