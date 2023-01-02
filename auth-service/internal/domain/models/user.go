package models

type User struct {
	ID        string
	Username  string //`validate:"string,max=256"`
	Firstname string
	Lastname  string
	Email     string //`validate:"email"`
	Phone     string
	Password  string
}

type ClaimsDTO struct {
	ID string
}
