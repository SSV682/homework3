package kafka

import (
	"order-service/internal/domain/dto"
	"time"
)

type ResponseCommand struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}

func (c *ResponseCommand) ToDTO() dto.OrderCommandDTO {
	return dto.OrderCommandDTO{
		OrderID: c.OrderID,
		Status:  dto.Status(c.Status),
	}
}

type RequestCommand struct {
	CommandType string `json:"command_type"`
	Order       Order  `json:"order"`
}

type Order struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	Status     string    `json:"status"`
}

func RequestCommandFromDTO(commandDTO dto.CommandDTO) RequestCommand {
	return RequestCommand{
		CommandType: string(commandDTO.CommandType),
		Order: Order{
			ID:         commandDTO.Order.ID,
			UserID:     commandDTO.Order.UserID,
			TotalPrice: commandDTO.Order.TotalPrice,
			CreatedAt:  commandDTO.Order.CreatedAt,
			Status:     commandDTO.Order.Status,
		},
	}
}
