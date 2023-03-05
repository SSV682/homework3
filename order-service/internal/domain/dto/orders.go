package dto

import (
	"encoding/json"
	"time"
)

type OrderRequestDTO struct {
	UserID         string
	IdempotencyKey string
	TotalPrice     float64
	DeliveryAt     time.Time
	Products       json.RawMessage
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
	DeliveryAt time.Time
	Products   json.RawMessage
	CreatedAt  time.Time
	Status     string
}

type OrdersDTO struct {
	Total   int
	Results []*OrderDTO
}
