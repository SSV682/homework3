package notification

type CreateProductRequest struct {
	Quantity int64  `json:"quantity"`
	Name     string `json:"name"`
}

type FillProductsRequest struct {
	Data []struct {
		ID       int64  `json:"id"`
		Quantity int64  `json:"quantity"`
		Name     string `json:"name"`
	} `json:"data"`
}

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseCreated struct {
	ID int64 `json:"id"`
}

type ResponseNotification struct {
	ID      int64  `json:"id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type ResponseNotifications struct {
	Total   int                     `json:"total"`
	Results []*ResponseNotification `json:"results"`
}
