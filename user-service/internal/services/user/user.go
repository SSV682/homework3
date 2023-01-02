package user

import (
	"context"
	"fmt"
	"user-service/internal/domain/models"
	"user-service/internal/provider"
)

type userService struct {
	sqlProv   provider.SqlUserProvider
	tokenProv provider.TokenProvider
}

func NewUserService(s provider.SqlUserProvider, t provider.TokenProvider) *userService {
	return &userService{
		sqlProv:   s,
		tokenProv: t,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) (string, error) {
	i, err := s.sqlProv.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return i, nil
}

func (s *userService) GetUser(ctx context.Context, token string) (models.User, error) {
	id, err := s.getIDFromClaims(token)
	if err != nil {
		return models.User{}, fmt.Errorf("get user doesnt have id: %v", err)
	}

	user, err := s.sqlProv.GetUser(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by id %d: %v", id, err)
	}
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, token string) error {
	id, err := s.getIDFromClaims(token)
	if err != nil {
		return fmt.Errorf("get user : %v", err)
	}

	if err = s.sqlProv.DeleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUser(ctx context.Context, token string, user *models.User) error {
	id, err := s.getIDFromClaims(token)
	if err != nil {
		return fmt.Errorf("get user : %v", err)
	}

	if err = s.sqlProv.UpdateUser(ctx, id, user); err != nil {
		return err
	}
	return nil
}

func (s *userService) getIDFromClaims(token string) (string, error) {
	//claims, err := s.tokenProv.ParseToken(token)
	//if err != nil {
	//	return "", fmt.Errorf("get id from claims: %v", err)
	//}
	//if ID, found := claims["id_user"]; found {
	//	return ID.(string), nil
	//} else {
	return "", fmt.Errorf("failed cast user_id")
	//}
}
