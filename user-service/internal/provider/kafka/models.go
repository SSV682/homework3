package kafka

import (
	domain "user-service/internal/domain/models"
)

type Command struct {
	UserID string `json:"user_id"`
	Mail   string `json:"mail"`
}

func NewCommand(command domain.User) Command {
	return Command{
		UserID: command.ID,
		Mail:   command.Email,
	}
}
