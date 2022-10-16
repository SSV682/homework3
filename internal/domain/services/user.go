package services

import (
	"context"
	"homework2/internal/models"
)

type UserService struct {
}

func (us *UserService) CreateUser(ctx context.Context, username, firstname, lastname, email, phone string, id int64) {
	user := &models.User{
		Id:        id,
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Phone:     phone,
	}
}
