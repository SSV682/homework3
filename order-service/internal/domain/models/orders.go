package domain

import (
	"order-service/internal/domain/dto"
	"time"
)

type StatusEnum string

type IntermediateOrderFunc func(o *Order) (bool, error)

type Order struct {
	id         int64
	userID     string
	totalPrice float64
	createdAt  time.Time
	status     Status
}

func (o *Order) ID() int64 {
	return o.id
}

func (o *Order) UserID() string {
	return o.userID
}

func (o *Order) TotalPrice() float64 {
	return o.totalPrice
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) SetStatus(status Status) {
	o.status = status
}

func NewOrderFromDTO(dto *dto.OrderRequestDTO) *Order {
	return &Order{
		userID:     dto.UserID,
		totalPrice: dto.TotalPrice,
		createdAt:  time.Now(),
		status:     Created,
	}
}

func RestoreOrderFromDTO(dto *dto.OrderDTO) *Order {
	return &Order{
		id:         dto.ID,
		userID:     dto.UserID,
		totalPrice: dto.TotalPrice,
		createdAt:  dto.CreatedAt,
		status:     Status(dto.Status),
	}
}

func (o *Order) OrderToDTO() *dto.OrderDTO {
	return &dto.OrderDTO{
		ID:         o.ID(),
		UserID:     o.UserID(),
		TotalPrice: o.TotalPrice(),
		CreatedAt:  o.CreatedAt(),
		Status:     string(o.Status()),
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
