package services

import (
	"context"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"user-service/internal/domain/models"
)

type AuthService interface {
	CheckUser(ctx context.Context, token string) (*models.ClaimsDTO, error)
	LoginUser(ctx context.Context, username, password string) (*string, error)
	GetKeys() (jwk.Set, error)
}
