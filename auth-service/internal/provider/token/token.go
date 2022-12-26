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
	SecretKey  = "AlexKraken"
	PrivateKey = "-----BEGIN RSA PRIVATE KEY-----\n" +
		"MIIBOwIBAAJBAMlxrcVZp+r6OZLC0smL/eUobx9T5dotNu6jD3I2sj7alyKM6Ylv" +
		"KgvL0MgyfYlyl6Ly32XxtvfA1vwnxvdmHxkCAwEAAQJAPq7R/Mv2NWchjSp0fuTB" +
		"35HiaiQoLOjO5Bj3UHn2oxnB7Nb0zpHZmfGwvYh+nz3vC5B5t6y5NRPsKLF0JWsW" +
		"kQIhAPjLqh10XfdUGajqGlJh+9Hjzv53mPOlQWYr7yezEJlFAiEAz0b+sfPEe4rA" +
		"WE7bpaPpITD+Bt2FDXb2rofEq0g1ycUCIARW0P24hNcGcXgftRvQt6qedYK8pT9C" +
		"l5RnmcEwf06dAiEAkWHbXNd8taZRWN8ewmRgLQ6e7hPLsfEB/tJtiDGiwH0CIQCz" +
		"UY/9L/HG/X28D5oKLgySx5YI/jYcVoLhomx7GCxSZQ==" +
		"\n-----END RSA PRIVATE KEY-----"
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

func (t *tokenProvider) CreateToken(userID string) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(PrivateKey))
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	claims := NewClaims(userID)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	headers := jwtToken.Header
	headers["kid"] = KID
	jwtToken.Header = headers

	jwtSignedB64, err := jwtToken.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("doesnt signed token: %v", err)
	}
	return jwtSignedB64, nil
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
