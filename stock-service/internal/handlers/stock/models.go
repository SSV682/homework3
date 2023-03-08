package stock

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

type ResponseProduct struct {
	ID       int64  `json:"id"`
	Quantity int64  `json:"quantity"`
	Name     string `json:"name"`
}

type ResponseProducts struct {
	Total   int                `json:"total"`
	Results []*ResponseProduct `json:"results"`
}
