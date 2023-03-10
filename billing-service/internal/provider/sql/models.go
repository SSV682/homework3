package sql

import (
	domain "billing-service/internal/domain/models"
	"encoding/json"
	"github.com/google/uuid"
)

type AccountRow struct {
	UserID uuid.UUID `db:"id"`
	Amount float64   `db:"amount"`
}

func FromModel(account *domain.Account) *AccountRow {
	return &AccountRow{
		UserID: account.UserID,
		Amount: account.Amount,
	}
}

func (r *AccountRow) ToModel() *domain.Account {
	return &domain.Account{
		UserID: r.UserID,
		Amount: r.Amount,
	}
}

type CommandRow struct {
	ID      int64           `db:"id"`
	Topic   string          `db:"topic"`
	Message json.RawMessage `db:"topic"`
}

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
