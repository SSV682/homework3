package provider

import (
	"context"

	"user-service/internal/domain/models"
)

type SqlUserProvider interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUser(ctx context.Context, id int64) (models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, id int64, user *models.User) error
}

type TokenProvider interface {
	ParseToken(tokenString string) (models.Claims, error)
}
