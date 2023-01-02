package models

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	ID string
	jwt.RegisteredClaims
}

type KeyDTO struct {
	ID        string
	PublicKey string
}
