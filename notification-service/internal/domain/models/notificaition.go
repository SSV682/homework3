package domain

type Status string

const (
	Success           Status = "success"
	Canceled          Status = "canceled"
	Created           Status = "created"
	PaymentPending    Status = "payment_pending"
	PaymentRejecting  Status = "payment_rejecting"
	PaymentRejected   Status = "payment_rejected"
	PaymentApproved   Status = "payment_approved"
	StockPending      Status = "stock_pending"
	StockApproved     Status = "stock_approved"
	StockRejected     Status = "stock_rejected"
	StockRejecting    Status = "stock_rejecting"
	DeliveryPending   Status = "delivery_pending"
	DeliveryApproved  Status = "delivery_approved"
	DeliveryRejected  Status = "delivery_rejected"
	DeliveryRejecting Status = "delivery_rejecting"
	Canceling         Status = "canceling"
	ErroneousAttempt  Status = "erroneous_attempt"
)

type Order struct {
	ID     int64
	UserID string
	Status Status
}

type Notification struct {
	ID      int64
	UserID  string
	Message string
}
