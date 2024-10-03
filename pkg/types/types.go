package types

type Student struct {
	ID    string `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email"`
	Age   int    `json:"age" validate:"required"`
}

type Response struct {
	Success bool `json:"success"`
	Message int  `json:"message"`
	Data    any  `json:"data"`
}