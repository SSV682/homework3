package token

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	Exp    int64  `json:"exp"`
	IDUser string `json:"id_user"`
}

func NewClaims(id string) Claims {
	exp := time.Now().Add(24 * time.Hour).Unix()

	return Claims{
		Exp:    exp,
		IDUser: id,
	}

}
