package token

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Exp    int64  `json:"exp"`
	IDUser string `json:"id_user"`
}
