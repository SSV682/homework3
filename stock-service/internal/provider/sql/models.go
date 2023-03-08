package sql

import domain "stock-service/internal/domain/models"

type ProductRow struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Quantity int64  `db:"quantity"`
}

func FromModel(order *domain.Product) *ProductRow {
	return &ProductRow{
		ID:       order.ID,
		Name:     order.Name,
		Quantity: order.Quantity,
	}
}

func (r *ProductRow) ToModel() *domain.Product {
	return &domain.Product{
		ID:       r.ID,
		Quantity: r.Quantity,
		Name:     r.Name,
	}
}
