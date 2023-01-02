package middleware

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	log "github.com/sirupsen/logrus"
	"strings"
)

func AuthMiddleware(url string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			header := ctx.Request().Header.Get("Authorization")
			jwtString := strings.Split(header, "Bearer ")[1]

			c := jwk.NewCache(context.Background())
			c.Register(url)
			keySet, err := c.Get(context.Background(), url)
			if err != nil {
				return fmt.Errorf("failed to fetch remote JWK: %s", err)
			}

			log.Info(keySet)
			verifiedToken, err := jwt.Parse([]byte(jwtString), jwt.WithValidate(true), jwt.WithKeySet(keySet))
			if err != nil {
				return fmt.Errorf("failed to verify JWS: %s\n", err)
			}

			if userID, found := verifiedToken.Get("user_id"); found {
				log.Infof("userID found: %s", userID.(string))
				ctx.Set("userID", userID)
			}

			err = next(ctx)
			if err != nil {
				ctx.Error(err)
			}
			return err
		}

	}
}
