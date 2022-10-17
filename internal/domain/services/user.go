package services

import (
	"context"
	"homework2/internal/domain/models"
	"time"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id int64) (models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, id int64, user *models.User) error
}

type UserService struct {
	userStorage    models.UserStorageInterface
	contextTimeout time.Duration
}

func NewUserService(u models.UserStorageInterface, timeout time.Duration) UserServiceInterface {
	return &UserService{
		userStorage:    u,
		contextTimeout: timeout,
	}
}

func (us *UserService) CreateUser(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	err := us.userStorage.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil

}

func (us *UserService) GetUser(ctx context.Context, id int64) (user models.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	user, err = us.userStorage.GetUser(ctx, id)
	if err != nil {
		return
	}
	return
}

func (us *UserService) DeleteUser(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	err := us.userStorage.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) UpdateUser(ctx context.Context, id int64, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()
	err := us.userStorage.UpdateUser(ctx, id, user)
	if err != nil {
		return err
	}
	return nil
}

//func (us *userService) UpdateUser()
