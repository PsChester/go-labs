package transport

import (
	"encoding/json"
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

type ShowOrdersRequestBody struct {
	UserId int `json:"user_id"`
}

type ShowOrderInfoRequestBody struct {
	OrderId string `json:"order_id"`
}

type CancelOrderRequestBody struct {
	OrderId string `json:"order_id"`
}

func Router(orderService model.OrderServiceInterface) http.Handler {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/").Subrouter()
	subRouter.HandleFunc("/cancel_order", CancelOrder(orderService)).Methods(http.MethodPost)
	subRouter.HandleFunc("/get_order", ShowOrderInfo(orderService)).Methods(http.MethodPost)
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

func ShowOrderInfo(orderService model.OrderServiceInterface) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		defer request.Body.Close()

		var message ShowOrderInfoRequestBody
		err = json.Unmarshal(body, &message)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		order, err := orderService.GetOrder(message.OrderId)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}

		jsonAnswer, err := json.Marshal(order)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}

		if _, err = io.WriteString(responseWriter, string(jsonAnswer)); err != nil {
			log.WithField("getOrderInfo", "failed")
		}
	}
}

func ShowOrders(orderService model.OrderServiceInterface) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		defer request.Body.Close()

		var message ShowOrdersRequestBody
		err = json.Unmarshal(body, &message)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		orders, err := orderService.GetAllOrders(message.UserId)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}

		jsonAnswer, err := json.Marshal(orders)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}

		if _, err = io.WriteString(responseWriter, string(jsonAnswer)); err != nil {
			log.WithField("getOrders", "failed") //TODO:question почему не отображается в логах?
		}
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

		err = orderService.CreateOrder(message.UserId, &message.ProductIds)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}
}

func CancelOrder(orderService model.OrderServiceInterface) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		defer request.Body.Close()

		var message CancelOrderRequestBody
		err = json.Unmarshal(body, &message)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		//TODO: проверка прав пользователя на отмену заказа

		err = orderService.CancelOrder(message.OrderId)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}
}
