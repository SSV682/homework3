package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

type tokenProvider struct {
	verifyKey *rsa.PublicKey
	secretKey string
}

const (
	SecretKey = "AlexKraken"
	PublicKey = "-----BEGIN PUBLIC KEY-----\n" +
		"MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAMlxrcVZp+r6OZLC0smL/eUobx9T5dotNu6jD3I2sj7alyKM6YlvKgvL0MgyfYlyl6Ly32XxtvfA1vwnxvdmHxkCAwEAAQ==" +
		"\n-----END PUBLIC KEY-----"
	KID = "zXew0UJ1h6Q4CCcd_9wxMzvcp5cEBifH0KWrCz2Kyxc"
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

func (t *tokenProvider) ParseToken(tokenString string) (jwt.MapClaims, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(PublicKey))
	if err != nil {
		return nil, fmt.Errorf("validate parse key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	log.Info(token.Claims)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf(ErrorInvalidToken, err)
	}
}
