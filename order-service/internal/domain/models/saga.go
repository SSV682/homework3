package domain

type Status string

const (
	Success          Status = "success"
	Canceled         Status = "canceled"
	Created          Status = "created"
	PaymentPending   Status = "payment_pending"
	PaymentRejecting Status = "payment_rejecting"
	PaymentRejected  Status = "payment_rejected"
	PaymentApproved  Status = "payment_approved"
	StockPending     Status = "stock_pending"
	StockApproved    Status = "stock_approved"
)

type CreateOrderSagaState struct {
	OrderID int64
	Status  Status
}
