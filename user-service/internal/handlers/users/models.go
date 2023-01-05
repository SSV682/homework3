package users

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseCreated struct {
	ID string `json:"id"`
}
