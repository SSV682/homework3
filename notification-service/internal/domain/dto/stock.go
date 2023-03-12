package dto

type ProductRequestDTO struct {
	Quantity int64
	Name     string
}

type FillRequestDTO struct {
	Data []FillProductRequestDTO
}

type FillProductRequestDTO struct {
	ID       int64
	Quantity int64
	Name     string
}

type FilterProductDTO struct {
	Quantity *int64
	Limit    uint64
	Offset   uint64
}
