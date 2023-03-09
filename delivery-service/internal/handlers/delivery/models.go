package delivery

import "encoding/json"

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseDelivery struct {
	ID           int64           `json:"id"`
	OrderId      int64           `json:"order_id"`
	OrderContent json.RawMessage `json:"order_content"`
	Address      json.RawMessage `json:"address"`
	DeliveryAt   int64           `json:"delivery_at"`
}

type ResponseDeliveries struct {
	Total   int                 `json:"total"`
	Results []*ResponseDelivery `json:"results"`
}
