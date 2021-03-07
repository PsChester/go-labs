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
	CreatedDate int       `json:"created_date"`
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
	panic("implement me")
	//1. Запросить все заказы из таблицы order с соответсвующим user_id
	//2. В цикле для каждого заказа запросить из бд список продуктов
	//3. Вернуть массив
}

func (orderService *OrderService) CancelOrder(orderId string) error {
	panic("implement me")
}
func (orderService *OrderService) GetOrder(orderId string) (Order, error) {
	panic("implement me")
}

func (orderService *OrderService) UpdateOrder(orderId string, productIds *[]int) error {
	panic("implement me")
}
