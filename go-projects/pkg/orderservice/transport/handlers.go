package transport

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type Order struct {
	Id    string `json:"id"` //TODO:Question почему с заглавно буквы?
	Price int    `json:"price"`
}

func Router() http.Handler {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/").Subrouter()
	subRouter.HandleFunc("/order/{id}", getOrder).Methods(http.MethodGet)
	return logMiddleware(router)
}

func logMiddleware(httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		log.WithFields(log.Fields{
			"method":     request.Method,
			"url":        request.URL,
			"remoteAddr": request.RemoteAddr,  ////TODO:Question что это?
			"userAgent":  request.UserAgent(), //информация о названии и версии приложения (браузера), операционную систему компьютера и язык..
			"time":       time.Now(),
		}).Info("got a new request")
		httpHandler.ServeHTTP(responseWriter, request)
	})
}

//Есть ли опциональный возврат? Order|nil
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
