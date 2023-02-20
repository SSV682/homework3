package dto

import "time"

type Status string

const (
	Success         Status = "success"
	Created         Status = "created"
	Canceling       Status = "canceling"
	PaymentPending  Status = "payment_pending"
	PaymentApproved Status = "payment_approved"
	StockPending    Status = "stock_pending"
	StockApproved   Status = "stock_approved"
	PaymentRejected Status = "payment_rejected"
	StockRejected   Status = "stock_rejected"
)

type CommandType string

const (
	Approve CommandType = "approve"
	Reject  CommandType = "reject"
)

type OrderRequestDTO struct {
	UserID         string
	IdempotencyKey string
	TotalPrice     float64
	CreatedAt      time.Time
}

type FilterOrderDTO struct {
	UserID string
	Limit  uint64
	Offset uint64
}

type OrderDTO struct {
	ID         int64
	UserID     string
	TotalPrice float64
	CreatedAt  time.Time
	Status     string
}

type CommandDTO struct {
	CommandType CommandType
	Order       OrderDTO
}

type OrdersDTO struct {
	Total   int
	Results []*OrderDTO
}

type OrderCommandDTO struct {
	OrderID int64
	Status  Status
}
