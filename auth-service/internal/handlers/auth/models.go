package auth

type ResponseError struct {
	Message string `json:"message"`
}

type User struct {
	Username string `param:"username" query:"username" json:"username" validate:"required"`
	Password string `param:"password" query:"password" json:"password" validate:"required,password"`
}