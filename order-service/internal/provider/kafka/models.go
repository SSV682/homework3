package kafka

import (
	"encoding/json"
	domain "order-service/internal/domain/models"
	"time"
)

type ResponseCommand struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}

func (c *ResponseCommand) ToDTO() domain.OrderCommand {
	return domain.OrderCommand{
		OrderID: c.OrderID,
		Status:  domain.Status(c.Status),
	}
}

type RequestCommand struct {
	CommandType string `json:"command_type"`
	Order       Order  `json:"order"`
}

type Order struct {
	ID         int64           `json:"id"`
	UserID     string          `json:"user_id"`
	TotalPrice float64         `json:"total_price"`
	Products   json.RawMessage `json:"products"`
	CreatedAt  time.Time       `json:"created_at"`
	DeliveryAt time.Time       `json:"delivery_at"`
	Status     string          `json:"status"`
}

func RequestCommandFromDTO(command domain.Command) RequestCommand {
	return RequestCommand{
		CommandType: string(command.CommandType),
		Order: Order{
			ID:         command.Order.ID,
			UserID:     command.Order.UserID,
			TotalPrice: command.Order.TotalPrice,
			Products:   command.Order.Products,
			CreatedAt:  command.Order.CreatedAt,
			DeliveryAt: command.Order.DeliveryAt,
			Status:     string(command.Order.Status),
		},
	}
}
