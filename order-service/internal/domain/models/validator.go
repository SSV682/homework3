package domain

type Validator interface {
	Struct(s interface{}) error
}
