package models

import "context"

type User struct {
	Id        int64
	Username  string //`validate:"string,max=256"`
	Firstname string
	Lastname  string
	Email     string //`validate:"email"`
	Phone     string
}

type UserStorageInterface interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int64) (User, error)
	DeleteUser(ctx context.Context, id int64) error
	//UpdateUser(ctx context.Context, id int64, user *User) error
}
