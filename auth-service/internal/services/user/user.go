package user

import (
	"auth-service/internal/domain/errors"
	"auth-service/internal/domain/models"
	"auth-service/internal/provider"
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/v2/jwk"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	sqlProv   provider.UserProvider
	tokenProv provider.TokenProvider
}

func NewUserService(s provider.UserProvider, t provider.TokenProvider) *userService {
	return &userService{
		sqlProv:   s,
		tokenProv: t,
	}
}

func (s *userService) LoginUser(ctx context.Context, username, password string) (*string, error) {
	user, err := s.sqlProv.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf(ErrorCheckUser, err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
		token, err := s.tokenProv.CreateToken(user.ID)
		if err != nil {
			return nil, fmt.Errorf(ErrorCreateToken, err)
		}

		return &token, err
	}

	return nil, fmt.Errorf(ErrorWrongPass, err)
}

func (s *userService) CheckUser(ctx context.Context, payload string) (*models.ClaimsDTO, error) {
	token, err := s.tokenProv.ParseToken(payload)
	if err != nil {
		log.Errorf(ErrorParseToken, err)
		return nil, errors.ErrFailedToken
	}

	userID, found := token.Get("id_user")
	if !found {
		log.Errorf("failed cast user_id")
		return nil, errors.ErrFailedToken
	}

	if _, err := s.sqlProv.GetUserByID(ctx, userID.(string)); err != nil {
		log.Errorf(ErrorCheckUser, err)
		return nil, errors.ErrInternalError
	}

	return &models.ClaimsDTO{ID: userID.(string)}, nil
}

func (s *userService) GetKeys() (jwk.Set, error) {
	set, err := s.tokenProv.GetKeys()
	if err != nil {
		log.Printf("failed to marshal key set into JSON: %s", err)

	}

	return set, nil
}
