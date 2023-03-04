package sql

import (
	"order-service/internal/domain/dto"
	domain "order-service/internal/domain/models"
	"time"
)

type OrderRow struct {
	ID         int64     `db:"id"`
	UserID     string    `db:"user_id"`
	TotalPrice float64   `db:"total_price"`
	CreatedAt  time.Time `db:"created_at"`
	Status     string    `db:"status"`
}

func (r *OrderRow) ToDTO() *dto.OrderDTO {
	return &dto.OrderDTO{
		ID:         r.ID,
		UserID:     r.UserID,
		TotalPrice: r.TotalPrice,
		CreatedAt:  r.CreatedAt,
		Status:     r.Status,
	}
}

func FromModel(order *domain.Order) *OrderRow {
	return &OrderRow{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		Status:     string(order.Status),
	}
}
