package transport

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)
import log "github.com/sirupsen/logrus"

func Router() http.Handler {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	return logMiddleware(router)
}

func logMiddleware(httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		log.WithFields(log.Fields{
			"method": request.Method,
			"url": request.URL,
			"remoteAddr": request.RemoteAddr,  //Что это?
			"userAgent": request.UserAgent(), //информация о названии и версии приложения (браузера), операционную систему компьютера и язык..
			"time": time.Now(),
		}).Info("got a new request")
		httpHandler.ServeHTTP(responseWriter, request)
	})
}

func helloWorld(responseWriter http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(responseWriter, "Hello world")
}