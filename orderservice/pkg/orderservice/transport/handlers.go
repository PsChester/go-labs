package transport

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io" //TODO:Question можно ли сделать 1 иморт io??
	"net/http"
	"time"
)

type Server struct {
	Database *sql.DB
}

// Order Answer: название структур и их полей с заглавной буквы, чтобы их можно было использовать при экспорте.
type Order struct {
	Id    string `json:"id"`
	Price int    `json:"price"`
}

func Router(server *Server) http.Handler {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/").Subrouter()
	subRouter.HandleFunc("/orders", server.showOrders).Methods(http.MethodGet)
	subRouter.HandleFunc("/order_creating", server.createOrder).Methods(http.MethodPost)
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

//TODO:Question Есть ли опциональный возврат? Order|nil
func (server *Server) createOrder(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	var message Order
	err = json.Unmarshal(body, &message)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	//TODO: проверка на пустоту body

	orderId := uuid.New().String()
	//TODO: вынести названия в константы, цену брать из body
	query := "INSERT INTO orderservice.order (order_id, price) VALUES (?, ?)"
	result, err := server.Database.Exec(query, orderId, 100)

	if err != nil {
		log.WithField("Database.Exec", "No added")
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonAnswer, err := json.Marshal(result.LastInsertId)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = io.WriteString(responseWriter, string(jsonAnswer)); err != nil {
		log.WithField("err", err).Error("write response err")
	}
}

func (server *Server) showOrders(responseWriter http.ResponseWriter, _ *http.Request) {
	query := "SELECT * FROM orderservice.order"
	rows, err := server.Database.Query(query)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	orders := make([]Order, 0)
	for rows.Next() {
		order := Order{}
		err = rows.Scan(&order.Id, &order.Price)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	jsonAnswer, err := json.Marshal(orders)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = io.WriteString(responseWriter, string(jsonAnswer)); err != nil {
		log.WithField("error", err).Error("write response err")
	}
}
