package domain

import (
	"stock-service/internal/domain/dto"
)

type Status string

const (
	StockApproved Status = "stock_approved"
	StockRejected Status = "stock_rejected"
)

type IntermediateProductFunc func(o *Product) (bool, error)

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

type Product struct {
	ID       int64
	Quantity int64
	Name     string
}

func RestoreProductFromRequest(dto dto.ProductRequestDTO) Product {
	return Product{
		Quantity: dto.Quantity,
		Name:     dto.Name,
	}
}

type Command struct {
	OrderID int64
	Status  Status
}

type ResponseCommand struct {
	Topic   string
	Command Command
}

type Order struct {
	ID       int64
	Products []Product
}

type RequestCommand struct {
	CommandType CommandType
	Order       Order
}
