package user

import (
	"context"
	"fmt"
	"user-service/internal/domain/models"
	"user-service/internal/provider"
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

	if user.Password == password {
		token, err := s.tokenProv.CreateToken(user.Id)
		if err != nil {
			return nil, fmt.Errorf(ErrorCreateToken, err)
		}

		return &token, err
	}

	return nil, fmt.Errorf(ErrorWrongPass, err)
}

func (s *userService) CheckUser(ctx context.Context, token string) (bool, error) {
	claims, err := s.tokenProv.ParseToken(token)
	if err != nil {
		return false, fmt.Errorf(ErrorParseToken, err)
	}

	if _, err := s.sqlProv.GetUserByID(ctx, claims.ID); err != nil {
		return false, fmt.Errorf(ErrorCheckUserByID, err)
	}

	//if claims.Expire.Unix() < time.Now().Unix() {
	//	return false, nil
	//}
	return true, nil
}

func (s *userService) SignUpUser(ctx context.Context, user *models.User) (string, error) {
	i, err := s.sqlProv.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return i, nil
}
