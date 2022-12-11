package token

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Exp    int64 `json:"exp"`
	IDUser int64 `json:"id_user"`
}
