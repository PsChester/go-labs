package transport

import (
//	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"io"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Kitty struct {
	Name string `json:"Name"` //TODO:Question так мы указываем конкретный вид строки??
}

func Router() http.Handler {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/").Subrouter()
	subRouter.HandleFunc("/kitty", getKitty).Methods(http.MethodGet)
	return logMiddleware(router)
}

func logMiddleware(httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		log.WithFields(log.Fields{
			"method": request.Method,
			"url": request.URL,
			"remoteAddr": request.RemoteAddr,  ////TODO:Question что это?
			"userAgent": request.UserAgent(), //информация о названии и версии приложения (браузера), операционную систему компьютера и язык..
			"time": time.Now(),
		}).Info("got a new request")
		httpHandler.ServeHTTP(responseWriter, request)
	})
}

func getKitty(responseWriter http.ResponseWriter, _ *http.Request) {
	cat := Kitty{"Кот"}
	jsonAnswer, error := json.Marshal(cat)
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