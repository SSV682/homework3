package provider

import (
	"context"

	"user-service/internal/domain/models"
)

type UserProvider interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id int64) (models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, id int64, user *models.User) error
}
