package transport

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil" //TODO:Question можно ли сделать 1 иморт io??
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
	subRouter.HandleFunc("/order/{id}", getOrder).Methods(http.MethodGet)
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
func getOrderById(id string) (Order, error) {
	orders := []Order{
		{
			Id:    "11",
			Price: 100,
		},
		{
			Id:    "12",
			Price: 200,
		},
	}

	for _, order := range orders {
		if order.Id == id {
			return order, nil
		}
	}

	return Order{}, errors.New("Order don't found")
}

func getOrder(responseWriter http.ResponseWriter, request *http.Request) {
	variables := mux.Vars(request)
	id := variables["id"]
	order, error := getOrderById(id)
	if error != nil {
		http.Error(responseWriter, error.Error(), http.StatusInternalServerError)
		return
	}

	jsonAnswer, error := json.Marshal(order)
	if error != nil {
		http.Error(responseWriter, error.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	responseWriter.WriteHeader(http.StatusOK)
	//Переменные объявленные внутри if коротким образом, также доступны внутри else блоков
	if _, error = io.WriteString(responseWriter, string(jsonAnswer)); error != nil {
		log.WithField("error", error).Error("write response error")
	}
}

func (server *Server) createOrder(responseWriter http.ResponseWriter, request *http.Request) {
	body, error := ioutil.ReadAll(request.Body)
	if error != nil {
		http.Error(responseWriter, error.Error(), http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	var message Order
	error = json.Unmarshal(body, &message)
	if error != nil {
		http.Error(responseWriter, error.Error(), http.StatusInternalServerError)
		return
	}
	//TODO: проверка на пустоту body

	orderId := uuid.New().String()
	//TODO: вынести названия в константы, цену брать из body
	query := "INSERT INTO orderservice.order (order_id, price) VALUES (?, ?)"
	result, error := server.Database.Exec(query, orderId, 100)

	if error != nil {
		log.WithField("Database.Exec", "No added")
		http.Error(responseWriter, error.Error(), http.StatusInternalServerError)
		return
	}

	jsonAnswer, error := json.Marshal(result.LastInsertId)
	if error != nil {
		http.Error(responseWriter, error.Error(), http.StatusInternalServerError)
		return
	}

	log.WithField("Check", "5")

	if _, error = io.WriteString(responseWriter, string(jsonAnswer)); error != nil {
		log.WithField("error", error).Error("write response error")
	}

	log.WithField("Check", "6")
}
