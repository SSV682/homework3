package billing

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseAccount struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

type RequestFillAccount struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

type RequestCreateAccount struct {
	ID string `json:"id"`
}
