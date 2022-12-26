package models

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	ID string
	jwt.RegisteredClaims
}
