package domain

import (
	"order-service/internal/domain/dto"
	"time"
)

type StatusEnum string

type IntermediateOrderFunc func(o *Order) (bool, error)

type CommandType string

const (
	Approve CommandType = "approve"
	Reject  CommandType = "reject"
)

type Order struct {
	ID         int64
	UserID     string
	TotalPrice float64
	CreatedAt  time.Time
	Status     Status
}

func (o *Order) SetStatus(status Status) {
	o.Status = status
}

func NewOrderFromDTO(dto *dto.OrderRequestDTO) *Order {
	return &Order{
		UserID:     dto.UserID,
		TotalPrice: dto.TotalPrice,
		CreatedAt:  time.Now(),
		Status:     Created,
	}
}

func RestoreOrderFromDTO(dto *dto.OrderDTO) *Order {
	return &Order{
		ID:         dto.ID,
		UserID:     dto.UserID,
		TotalPrice: dto.TotalPrice,
		CreatedAt:  dto.CreatedAt,
		Status:     Status(dto.Status),
	}
}

func (o *Order) OrderToDTO() *dto.OrderDTO {
	return &dto.OrderDTO{
		ID:         o.ID,
		UserID:     o.UserID,
		TotalPrice: o.TotalPrice,
		CreatedAt:  o.CreatedAt,
		Status:     string(o.Status),
	}
}

type Orders struct {
	total   int
	results []*Order
}

func (o *Orders) Total() int {
	return o.total
}

func (o *Orders) Results() []*Order {
	return o.results
}

func (o *Orders) OrdersToDTO() *dto.OrdersDTO {
	results := make([]*dto.OrderDTO, 0, len(o.Results()))

	for _, v := range o.Results() {
		results = append(results, v.OrderToDTO())
	}

	return &dto.OrdersDTO{
		Total:   o.Total(),
		Results: results,
	}
}

func NewOrdersFromSlice(res []*Order) *Orders {
	return &Orders{
		total:   len(res),
		results: res,
	}
}

type Command struct {
	Topic       string
	CommandType CommandType
	Order       *Order
}

type Message struct {
	Topic string
	Order Order
}

type OrderCommand struct {
	OrderID int64
	Status  Status
}
