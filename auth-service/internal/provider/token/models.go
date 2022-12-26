package token

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	IDUser string `json:"id_user"`
}

func NewClaims(userID string) Claims {
	exp := jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)}
	now := jwt.NumericDate{Time: time.Now()}
	uuidToken := uuid.New()

	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "http://userservice-authservice.userservice.svc.cluster.local",
			Subject:   "",
			Audience:  nil,
			ExpiresAt: &exp,
			NotBefore: &now,
			IssuedAt:  &now,
			ID:        uuidToken.String(),
		},
		IDUser: userID,
	}

}
