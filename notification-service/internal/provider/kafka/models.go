package kafka

import (
	domain "notification-service/internal/domain/models"
)

type Order struct {
	ID     int64  `json:"ID"`
	UserID string `json:"UserID"`
	Status string `json:"Status"`
}

func (c *Order) ToModel() domain.Order {
	return domain.Order{
		ID:     c.ID,
		UserID: c.UserID,
		Status: domain.Status(c.Status),
	}

}
