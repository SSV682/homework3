package orders

import (
	"encoding/json"
	"time"
)

type CreateOrderRequest struct {
	TotalPrice float64         `json:"total_price"`
	Products   json.RawMessage `json:"products"`
	DeliveryAt time.Time       `json:"delivery_at"`
}

type CancelingOrderRequest struct {
	Status float64 `json:"status"`
}

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseCreated struct {
	ID int64 `json:"id"`
}

type ResponseOrder struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	Status     string    `json:"status"`
}

type ResponseOrders struct {
	Total   int              `json:"total"`
	Results []*ResponseOrder `json:"results"`
}
