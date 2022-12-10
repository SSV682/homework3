package user

import (
	"context"
	"user-service/internal/domain/models"
	"user-service/internal/services"
)

type userService struct {
	sqlProv services.UserValueService
}

func NewUserService(s services.UserValueService) *userService {
	return &userService{
		sqlProv: s,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	err := s.sqlProv.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) GetUser(ctx context.Context, id int64) (user models.User, err error) {
	user, err = s.sqlProv.GetUser(ctx, id)
	if err != nil {
		return
	}
	return
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	err := s.sqlProv.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUser(ctx context.Context, id int64, user *models.User) error {
	err := s.sqlProv.UpdateUser(ctx, id, user)
	if err != nil {
		return err
	}
	return nil
}
