package provider

import (
	"context"
	"user-service/internal/domain/models"
)

type UserProvider interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
}

type TokenProvider interface {
	ParseToken(tokenString string) (models.Claims, error)
	CreateToken(userID int64) (string, error)
}
