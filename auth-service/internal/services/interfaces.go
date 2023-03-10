package services

import (
	"auth-service/internal/domain/models"
	"context"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type AuthService interface {
	CheckUser(ctx context.Context, token string) (*models.ClaimsDTO, error)
	LoginUser(ctx context.Context, username, password string) (*string, error)
	GetKeys() (jwk.Set, error)
}
