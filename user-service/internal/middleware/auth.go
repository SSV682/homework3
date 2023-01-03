package middleware

import (
	"encoding/base64"
	"fmt"
	"github.com/goccy/go-json"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomClaims struct {
	jwtv4.RegisteredClaims
	IDUser string `json:"id_user"`
}

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if payload := ctx.Request().Header.Get("x-jwt-token"); payload != "" {
				data, err := base64.StdEncoding.DecodeString(payload)
				if err != nil {
					return ctx.JSON(http.StatusUnauthorized, fmt.Errorf("couldnt decode payload: %s", err))
				}
				var claims CustomClaims

				err = json.Unmarshal(data, &claims)
				if err != nil {
					return ctx.JSON(http.StatusUnauthorized, fmt.Errorf("couldnt unmarshal payload: %s", err))
				}

				ctx.Set("userID", claims.IDUser)
				return next(ctx)
			}
			return ctx.JSON(http.StatusUnauthorized, fmt.Errorf("payload is empty"))
		}
	}
}
