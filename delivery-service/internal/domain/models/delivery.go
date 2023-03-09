package domain

import (
	"encoding/json"
	"time"
)

type Status string

const (
	DeliveryApproved Status = "delivery_approved"
	DeliveryRejected Status = "delivery_rejected"
)

type DeliveryEntry struct {
	ID           int64
	OrderID      int64
	OrderContent json.RawMessage
	Address      json.RawMessage
	Date         time.Time
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

type Order struct {
	ID           int64
	Address      json.RawMessage
	OrderContent json.RawMessage
	Date         time.Time
}

type RequestCommand struct {
	CommandType CommandType
	Order       Order
}
