package provider

import (
	"context"
	"user-service/internal/domain/models"
)

type UserProvider interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (string, error)
}

type TokenProvider interface {
	ParseToken(tokenString string) (models.Claims, error)
	CreateToken(userID string) (string, error)
}
