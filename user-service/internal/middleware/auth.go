package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	log "github.com/sirupsen/logrus"
	"strings"
)

const userIDFieldName = "id_user"

func AuthMiddleware(jwkURL string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if payload := ctx.Request().Header.Get("x-jwt-token"); payload != "" {
				data, err := base64.StdEncoding.DecodeString(payload)
				if err != nil {
					log.Fatal("error:", err)
				}
				log.Infof("payload: %v", payload)
				log.Infof("payload data: %v", data)
			}

			header := ctx.Request().Header.Get("Authorization")
			jwtString := strings.Split(header, "Bearer ")[1]

			c := jwk.NewCache(context.Background())
			c.Register(jwkURL)
			keySet, err := c.Get(context.Background(), jwkURL)
			if err != nil {
				return fmt.Errorf("failed to fetch remote JWK: %s", err)
			}

			log.Info(keySet)
			verifiedToken, err := jwt.Parse([]byte(jwtString), jwt.WithValidate(true), jwt.WithKeySet(keySet))
			if err != nil {
				return fmt.Errorf("failed to verify JWS: %s\n", err)
			}

			if userID, found := verifiedToken.Get(userIDFieldName); found {
				userValue := userID.(string)
				log.Infof("userID found: %s", userValue)
				ctx.Set("userID", userID.(string))
				//ctx.G
				return next(ctx)
			}
			log.Infof("userID not found: %v", verifiedToken)
			return fmt.Errorf("userID not found")
		}

	}
}
