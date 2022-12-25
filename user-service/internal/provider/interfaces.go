package provider

import (
	"context"

	"user-service/internal/domain/models"
)

type SqlUserProvider interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUser(ctx context.Context, id string) (models.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, user *models.User) error
}

type TokenProvider interface {
	ParseToken(tokenString string) (models.Claims, error)
}
