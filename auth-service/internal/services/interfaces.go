package services

import "context"

type AuthService interface {
	CheckUser(ctx context.Context, token string) (bool, error)
	LoginUser(ctx context.Context, username, password string) (*string, error)
}
