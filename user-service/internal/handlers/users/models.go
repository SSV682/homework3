package users

import "user-service/internal/domain/models"

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseCreated struct {
	ID string `json:"id"`
}

type Element struct {
	ID        string `json:"id"`
	Username  string `json:"username"` //`validate:"string,max=256"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"` //`validate:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

func (e *Element) ToModel() *models.User {
	return &models.User{
		ID:        e.ID,
		Username:  e.Username,
		Firstname: e.Firstname,
		Lastname:  e.Lastname,
		Email:     e.Email,
		Phone:     e.Phone,
		Password:  e.Password,
	}
}

func ToElement(m *models.User) Element {
	return Element{
		ID:        m.ID,
		Username:  m.Username,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
		Email:     m.Email,
		Phone:     m.Phone,
		Password:  m.Password,
	}
}
