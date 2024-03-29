package orders

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	IDUser string `json:"id_user"`
}

func getUserID(payload string) (string, error) {
	if payload != "" {
		data, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return "", fmt.Errorf("decode payload: %s", err)
		}

		var claims CustomClaims

		err = json.Unmarshal(data, &claims)
		if err != nil {
			return "", fmt.Errorf("unmarshal payload: %s", err)
		}

		if claims.IDUser != "" {
			return claims.IDUser, nil
		}
	}

	return "", fmt.Errorf("payload is empty")
}

func queryParamsToUInt64(value string, baseValue uint64) (uint64, error) {
	if value == "" {
		return baseValue, nil
	}

	limit, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return uint64(limit), nil
}
