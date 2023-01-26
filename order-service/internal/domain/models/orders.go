package domain

import (
	"time"
	"user-service/internal/domain/dto"
)

type StatusEnum string

func (e StatusEnum) String() string {
	return e.String()
}

const (
	StatusCreated         StatusEnum = "created"
	StatusPaymentAwaiting StatusEnum = "awaiting_payment"
	StatusPaymentReceived StatusEnum = "payment_received"
	StatusCompleted       StatusEnum = "completed"
	StatusCanceled        StatusEnum = "canceled"
	StatusFailed          StatusEnum = "failed"
)

type Order struct {
	id         int64
	userID     string
	totalPrice float64
	createdAt  time.Time
	status     StatusEnum
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

func (o *Order) Status() StatusEnum {
	return o.status
}

func NewOrderFromDTO(dto *dto.OrderRequestDTO) *Order {
	return &Order{
		userID:     dto.UserID,
		totalPrice: dto.TotalPrice,
		createdAt:  time.Now(),
		status:     StatusCreated,
	}
}

func (o *Order) OrderToDTO() *dto.OrderDTO {
	return &dto.OrderDTO{
		ID:         o.ID(),
		UserID:     o.UserID(),
		TotalPrice: o.TotalPrice(),
		CreatedAt:  o.CreatedAt(),
		Status:     o.Status().String(),
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
	results := make([]*Order, 0, len(res))
	for _, v := range res {
		results = append(results, &Order{
			id:         v.ID(),
			userID:     v.UserID(),
			totalPrice: v.TotalPrice(),
			createdAt:  v.CreatedAt(),
			status:     v.Status(),
		})
	}

	return &Orders{
		total:   len(results),
		results: results,
	}
}
