package domain

import (
	"github.com/google/uuid"
)

type Status string

const (
	PaymentRejected Status = "payment_rejected"
	PaymentApproved Status = "payment_approved"
)

type Account struct {
	UserID uuid.UUID
	Amount float64
}

type CommandType string

func ToCommandType(commandType string) CommandType {
	switch commandType {
	case "approve":
		return Approve
	case "reject":
		return Reject
	default:
		return Unknown
	}
}

const (
	Approve CommandType = "approve"
	Reject  CommandType = "reject"
	Unknown CommandType = "unknown"
)

type Command struct {
	OrderID int64
	Status  Status
}

type ResponseCommand struct {
	Topic   string
	Command Command
}

type ReadyResponseCommand struct {
	ID      int64
	Topic   string
	Command []byte
}

type Order struct {
	ID         int64
	TotalPrice float64
	UserID     uuid.UUID
}

type RequestCommand struct {
	CommandType CommandType
	Order       Order
}
