package main

import (
	"fmt"
	"net/http"
	"orderservice/transport"
	"os"
)
import log "github.com/sirupsen/logrus"

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)  //Как работает?
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	serverUrl := ":8000"
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")

	router := transport.Router() //http.Handler
	fmt.Println(http.ListenAndServe(":8000", router))
}
