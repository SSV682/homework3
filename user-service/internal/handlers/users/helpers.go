package users

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v4"
	"gopkg.in/go-playground/validator.v9"
	"user-service/internal/domain/models"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	IDUser string `json:"id_user"`
}

func getUserID(payload string) (string, error) {
	if payload != "" {
		data, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return "", fmt.Errorf("couldnt decode payload: %s", err)
		}

		var claims CustomClaims

		err = json.Unmarshal(data, &claims)
		if err != nil {
			return "", fmt.Errorf("couldnt unmarshal payload: %s", err)
		}

		if claims.IDUser != "" {
			return claims.IDUser, nil
		}
	}
	return "", fmt.Errorf("payload is empty")
}

func isRequestValid(u *models.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return false, err
	}
	return true, nil
}
