package models

type User struct {
	Id        int64
	Username  string //`validate:"string,max=256"`
	Firstname string
	Lastname  string
	Email     string //`validate:"email"`
	Phone     string
	Password  string
}
