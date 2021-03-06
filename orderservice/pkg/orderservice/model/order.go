package model

// Order Answer: название структур и их полей с заглавной буквы, чтобы их можно было использовать при экспорте.
type Order struct {
	Id    string `json:"id"`
	Price int    `json:"price"`
}
