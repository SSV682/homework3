package sql

import (
	domain "delivery-service/internal/domain/models"
	"encoding/json"
	"time"
)

type DeliveryRow struct {
	ID           int64           `db:"id"`
	OrderID      int64           `db:"order_id"`
	OrderContent json.RawMessage `db:"order_content"`
	Address      json.RawMessage `db:"address"`
	Date         time.Time       `db:"date"`
}

func FromModel(delivery *domain.DeliveryEntry) *DeliveryRow {
	return &DeliveryRow{
		ID:           delivery.ID,
		OrderID:      delivery.OrderID,
		OrderContent: delivery.OrderContent,
		Address:      delivery.Address,
		Date:         delivery.Date,
	}
}

func (r *DeliveryRow) ToModel() *domain.DeliveryEntry {
	return &domain.DeliveryEntry{
		ID:           r.ID,
		OrderID:      r.OrderID,
		OrderContent: r.OrderContent,
		Address:      r.Address,
		Date:         r.Date,
	}
}
