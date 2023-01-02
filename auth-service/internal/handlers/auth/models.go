package auth

type ResponseError struct {
	Message string `json:"message"`
}

type User struct {
	Username string `param:"username" query:"username" json:"username" validate:"required"`
	Password string `param:"password" query:"password" json:"password" validate:"required,password"`
}

type Token struct {
	Token string `param:"token" query:"token" json:"token" validate:"required"`
}

type Key struct {
	ID        string `json:"id"`
	PublicKey string `json:"public_key"`
}

type Keys struct {
	Keys []Key `json:"keys"`
}

type UserResponse struct {
	UserID string `json:"user_id"`
}
