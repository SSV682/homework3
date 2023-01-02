package provider

import (
	"context"
	//"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"user-service/internal/domain/models"
)

type UserProvider interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (string, error)
}

type TokenProvider interface {
	ParseToken(tokenString string) (jwt.Token, error)
	CreateToken(userID string) (string, error)
	GetKeys() (jwk.Set, error)
}
