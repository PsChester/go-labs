package model

import (
	"database/sql"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type OrderService struct {
	Database *sql.DB
}

// Order Answer: название структур и их полей с заглавной буквы, чтобы их можно было использовать при экспорте.
type Order struct {
	Id          string    `json:"id"`
	CreatedDate time.Time `json:"created_date"`
	Products    []Product `json:"products"`
}

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

//TODO:Question норм ли использовать Int для id-шников, которые в базе данных инкрементируются
type OrderServiceInterface interface {
	CreateOrder(userId int, productIds *[]int) error
	CancelOrder(orderId string) error
	GetOrder(orderId string) (Order, error)
	UpdateOrder(orderId string, productIds *[]int) error
	GetAllOrders(userId int) ([]Order, error)
}

//TODO:Question Есть ли опциональный возврат? Order|nil
func (orderService *OrderService) CreateOrder(userId int, productIds *[]int) error {
	//Тут должна быть проверка существования пользователя и продуктов

	orderId := uuid.New().String()
	query := "INSERT INTO orderservice.`order` (order_id, user_id, created_date) VALUES (?, ?, ?)"

	//TODO: Проверить мб нужно конвертировать time.Now() в timestamp
	_, err := orderService.Database.Exec(query, orderId, userId, time.Now())
	if err != nil {
		log.WithField("create_order", "failed")
		return err
	}

	for _, productId := range *productIds {
		query = "INSERT INTO orderservice.product_in_order (product_id, order_id) VALUES (?, ?)"
		_, err := orderService.Database.Query(query, productId, orderId)
		if err != nil {
			log.WithField("create_order", "failed")
			return err
		}
	}

	return nil
}

func (orderService *OrderService) GetAllOrders(userId int) ([]Order, error) {
	query := "SELECT order_id, created_date FROM orderservice.`order` WHERE user_id = ?"
	rows, err := orderService.Database.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]Order, 0)
	for rows.Next() {
		var order Order
		err = rows.Scan(&order.Id, &order.CreatedDate)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	for orderIndex, order := range orders {
		query = "SELECT product_id FROM orderservice.product_in_order WHERE order_id = ?"
		rows, err = orderService.Database.Query(query, order.Id)
		if err != nil {
			return nil, err
		}

		productIds := make([]string, 0)
		for rows.Next() {
			var productId string
			err = rows.Scan(&productId)
			if err != nil {
				return nil, err
			}
			productIds = append(productIds, productId)
		}

		orderProducts := make([]Product, 0)
		for _, productId := range productIds {
			query = "SELECT * FROM orderservice.product WHERE product_id = ?"
			rows, err = orderService.Database.Query(query, productId)
			if err != nil {
				return nil, err
			}

			for rows.Next() {
				var product Product
				err = rows.Scan(&product.Id, &product.Name, &product.Price)
				if err != nil {
					return nil, err
				}
				orderProducts = append(orderProducts, product)
			}
		}

		orders[orderIndex].Products = orderProducts
	}

	return orders, nil
}

func (orderService *OrderService) CancelOrder(orderId string) error {
	panic("implement me")
}
func (orderService *OrderService) GetOrder(orderId string) (Order, error) {
	query := "SELECT created_date FROM orderservice.`order` WHERE order_id = ?"
	var order Order
	err := orderService.Database.QueryRow(query, orderId).Scan(&order.CreatedDate)
	if err != nil {
		return Order{}, err
	}
	order.Id = orderId

	query = "SELECT product_id FROM orderservice.product_in_order WHERE order_id = ?"
	rows, err := orderService.Database.Query(query, order.Id)
	if err != nil {
		return Order{}, err
	}

	productIds := make([]string, 0)
	for rows.Next() {
		var productId string
		err = rows.Scan(&productId)
		if err != nil {
			return Order{}, err
		}
		productIds = append(productIds, productId)
	}

	order.Products = make([]Product, 0)
	for _, productId := range productIds {
		query = "SELECT * FROM orderservice.product WHERE product_id = ?"
		rows, err = orderService.Database.Query(query, productId)
		if err != nil {
			return Order{}, err
		}

		for rows.Next() {
			var product Product
			err = rows.Scan(&product.Id, &product.Name, &product.Price)
			if err != nil {
				return Order{}, err
			}
			order.Products = append(order.Products, product)
		}
	}

	return order, nil
}

func (orderService *OrderService) UpdateOrder(orderId string, productIds *[]int) error {
	panic("implement me")
}
