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

type OrderServiceInterface interface {
	CreateOrder(userId string, products []Product) error
	CancelOrder(orderId string)
	GetOrder(orderId string)
	UpdateOrder(orderId string, products []Product)
	GetAllOrders(userId string) ([]Order, error)
}

//TODO:Question Есть ли опциональный возврат? Order|nil
func (orderService *OrderService) CreateOrder(userId string, products []Product) error {
	orderId := uuid.New().String()
	createOrderQuery := "INSERT INTO orderservice.order (order_id, user_id, created_date) VALUES (?, ?, ?)"

	//TODO: Проверить мб нужно конвертировать time.Now() в timestamp
	_, err := orderService.Database.Exec(createOrderQuery, orderId, userId, time.Now())
	if err != nil {
		log.WithField("create_order", "failed")
		return err
	}

	//TODO: В цикле пополнить таблицу продуктами

	return nil
}

func (orderService *OrderService) GetAllOrders(userId string) ([]Order, error) {
	panic("implement me")
	//1. Запросить все заказы из таблицы order с соответсвующим user_id
	//2. В цикле для каждого заказа запросить из бд список продуктов
	//3. Вернуть массив
}

func (orderService *OrderService) CancelOrder(orderId string) {
	panic("implement me")
}
func (orderService *OrderService) GetOrder(orderId string) {
	panic("implement me")
}

func (orderService *OrderService) UpdateOrder(orderId string, products []Product) {
	panic("implement me")
}
