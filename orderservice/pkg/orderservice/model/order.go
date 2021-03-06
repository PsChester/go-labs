package model

// Order Answer: название структур и их полей с заглавной буквы, чтобы их можно было использовать при экспорте.
type Order struct {
	Id          string `json:"id"`
	CreatedDate int    `json:"created_date"`
}

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type OrderInterface interface {
	create(userId string)
	cancel(orderId string, userId string)
	get(orderId string, userId string)
	update(orderId string, userId string, products []Product)
}
