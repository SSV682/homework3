package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
	"user-service/internal/domain/models"
)

type tokenProvider struct {
	verifyKey *rsa.PublicKey
	secretKey string
}

const (
	SecretKey = "AlexKraken"
)

const (
	ErrorParseWithClaims = "parse with claims: %v"
	ErrorInvalidToken    = "invalid token: %v"
)

func NewJWTProvider() *tokenProvider {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(SecretKey))
	if err != nil {

	}
	return &tokenProvider{
		verifyKey: key,
		secretKey: SecretKey,
	}
}

func (t *tokenProvider) CreateToken(userID int64) (string, error) {
	claims := NewClaims(userID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(t.secretKey))
}

func (t *tokenProvider) ParseToken(tokenString string) (models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
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
