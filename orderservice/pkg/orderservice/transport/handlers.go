package transport

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"orderservice/pkg/orderservice/model"
	"time"
)

//TODO:Question можно ли считать данные из тела запроса, не создавая дополнительную структуру?
type CreateOrderRequestBody struct {
	UserId     int   `json:"user_id"`
	ProductIds []int `json:"product_ids"`
}

func Router(orderService model.OrderServiceInterface) http.Handler {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/").Subrouter()
	subRouter.HandleFunc("/get_orders", ShowOrders(orderService)).Methods(http.MethodPost)
	subRouter.HandleFunc("/create_order", CreateOrder(orderService)).Methods(http.MethodPost)
	return logMiddleware(router)
}

func logMiddleware(httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		log.WithFields(log.Fields{
			"method":     request.Method,
			"url":        request.URL,
			"remoteAddr": request.RemoteAddr,  ////TODO:Question что это?
			"userAgent":  request.UserAgent(), //Answer: информация о названии и версии приложения (браузера), операционную систему компьютера и язык..
			"time":       time.Now(),
		}).Info("got a new request")
		httpHandler.ServeHTTP(responseWriter, request)
	})
}

func ShowOrders(orderService model.OrderServiceInterface) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		//Получение user_id

		//orderService.GetAllOrders()
	}
}

func CreateOrder(orderService model.OrderServiceInterface) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		defer request.Body.Close()

		var message CreateOrderRequestBody
		err = json.Unmarshal(body, &message)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(message.UserId, message.ProductIds, body)

		err = orderService.CreateOrder(message.UserId, &message.ProductIds)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}
}
