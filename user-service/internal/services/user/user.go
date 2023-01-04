package user

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"user-service/internal/domain/models"
	"user-service/internal/provider"
)

type userService struct {
	sqlProv provider.SqlUserProvider
}

func NewUserService(s provider.SqlUserProvider) *userService {
	return &userService{
		sqlProv: s,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) (string, error) {
	password := user.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("password hashing failed:  %s", err)
	}
	user.Password = string(hashedPassword)

	i, err := s.sqlProv.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return i, nil
}

func (s *userService) GetUser(ctx context.Context, userID string) (models.User, error) {
	user, err := s.sqlProv.GetUser(ctx, userID)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by id %s: %v", userID, err)
	}
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	if err := s.sqlProv.DeleteUser(ctx, userID); err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUser(ctx context.Context, userID string, user *models.User) error {
	if err := s.sqlProv.UpdateUser(ctx, userID, user); err != nil {
		return err
	}
	return nil
}
