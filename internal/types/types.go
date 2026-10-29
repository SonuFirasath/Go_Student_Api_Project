package types

type Student struct {
	Id int64 `json:"id"`
	Name string `validate:"required"`
	Email string `validate:"required,email"`
	Age int `validate:"required,min=0"`
}