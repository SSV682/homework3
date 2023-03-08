package kafka

import (
	"errors"
	domain "stock-service/internal/domain/models"
	"time"
)

type Command struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}

func NewResponseCommand(command domain.Command) Command {
	return Command{
		OrderID: command.OrderID,
		Status:  string(command.Status),
	}
}

type Product struct {
	ID       int64  `json:"id"`
	Quantity int64  `json:"quantity"`
	Name     string `json:"name"`
}

type Order struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	Products   []Product `json:"products"`
	CreatedAt  time.Time `json:"created_at"`
	DeliveryAt time.Time `json:"delivery_at"`
	Status     string    `json:"status"`
}

type RequestCommand struct {
	CommandType string `json:"command_type"`
	Order       Order  `json:"order"`
}

func (c *RequestCommand) ToModel() (domain.RequestCommand, error) {
	ct := domain.ToCommandType(c.CommandType)
	if ct == domain.Unknown {
		return domain.RequestCommand{}, errors.New("bad command")
	}

	products := make([]domain.Product, 0, len(c.Order.Products))
	for _, v := range c.Order.Products {
		products = append(products, domain.Product{
			ID:       v.ID,
			Quantity: v.Quantity,
			Name:     v.Name,
		})
	}

	return domain.RequestCommand{
		CommandType: ct,
		Order: domain.Order{
			ID:       c.Order.ID,
			Products: products,
		},
	}, nil
}
