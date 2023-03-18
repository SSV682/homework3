package kafka

import (
	domain "billing-service/internal/domain/models"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
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

type Order struct {
	ID         int64           `json:"id"`
	UserID     uuid.UUID       `json:"user_id"`
	TotalPrice float64         `json:"total_price"`
	Products   json.RawMessage `json:"products"`
	Address    json.RawMessage `json:"address"`
	CreatedAt  time.Time       `json:"created_at"`
	DeliveryAt time.Time       `json:"delivery_at"`
	Status     string          `json:"status"`
}

type BillingRequestCommand struct {
	CommandType string `json:"command_type"`
	Order       Order  `json:"order"`
}

func (c *BillingRequestCommand) ToModel() (domain.RequestCommand, error) {
	ct := domain.ToCommandType(c.CommandType)
	if ct == domain.Unknown {
		return domain.RequestCommand{}, errors.New("bad command")
	}

	return domain.RequestCommand{
		CommandType: ct,
		Order: domain.Order{
			ID:         c.Order.ID,
			TotalPrice: c.Order.TotalPrice,
			UserID:     c.Order.UserID,
		},
	}, nil
}

type UserRequestCommand struct {
	UserID string `json:"user_id"`
	Mail   string `json:"mail"`
}

func (c *UserRequestCommand) ToModel() (domain.Account, error) {
	userID, err := uuid.FromBytes([]byte(c.UserID))
	if err != nil {
		return domain.Account{}, err
	}

	return domain.Account{
		UserID: userID,
		Amount: 0,
	}, nil
}
