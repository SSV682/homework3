package token

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
	"user-service/internal/domain/models"
	"user-service/internal/provider"
)

type tokenProvider struct {
}

const (
	SecretKey = "AlexKraken"
)

const (
	ErrorParseWithClaims = "parse with claims: %v"
	ErrorInvalidToken    = "invalid token: %v"
)

var Prov provider.TokenProvider = &tokenProvider{}

func (t *tokenProvider) CreateToken(userID int64) (string, error) {
	claims := NewClaims(userID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(SecretKey))
}

func (t *tokenProvider) ParseToken(tokenString string) (models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return models.Claims{}, fmt.Errorf(ErrorParseWithClaims, err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return models.Claims{
			ID:     claims.IDUser,
			Expire: time.Unix(claims.Exp, 0),
		}, nil
	} else {
		return models.Claims{}, fmt.Errorf(ErrorInvalidToken, err)
	}

}
