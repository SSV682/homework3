package services

import (
	"context"
	"user-service/internal/domain/models"
)

type AuthService interface {
	CheckUser(ctx context.Context, token string) (bool, error)
	SignUpUser(ctx context.Context, user *models.User) (string, error)
	LoginUser(ctx context.Context, username, password string) (*string, error)
}
