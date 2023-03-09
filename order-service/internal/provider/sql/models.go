package sql

import (
	"encoding/json"
	"order-service/internal/domain/dto"
	domain "order-service/internal/domain/models"
	"time"
)

type OrderRow struct {
	ID         int64           `db:"id"`
	UserID     string          `db:"user_id"`
	TotalPrice float64         `db:"total_price"`
	CreatedAt  time.Time       `db:"created_at"`
	DeliveryAt time.Time       `db:"delivery_at"`
	Address    json.RawMessage `db:"address"`
	Products   json.RawMessage `db:"products"`
	Status     string          `db:"status"`
}

func (r *OrderRow) ToDTO() *dto.OrderDTO {
	return &dto.OrderDTO{
		ID:         r.ID,
		UserID:     r.UserID,
		TotalPrice: r.TotalPrice,
		CreatedAt:  r.CreatedAt,
		DeliveryAt: r.DeliveryAt,
		Address:    r.Address,
		Products:   r.Products,
		Status:     r.Status,
	}
}

func FromModel(order *domain.Order) *OrderRow {
	return &OrderRow{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		DeliveryAt: order.DeliveryAt,
		Products:   order.Products,
		Address:    order.Address,
		Status:     string(order.Status),
	}
}
