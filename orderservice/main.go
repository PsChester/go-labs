package main

import (
//	"fmt"
	"net/http"
	"orderservice/transport"
	"os"
	"os/signal"
	"syscall"
	"context"
)
import log "github.com/sirupsen/logrus"

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)  ////TODO:Question Как работает?
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	serverUrl := ":8000"
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")
	killSignalChan := getKillSignalChan()
	server := startServer(serverUrl)

	waitForKillSignal(killSignalChan)
	server.Shutdown(context.Background())
}

func startServer(serverUrl string) *http.Server {
	router := transport.Router() //http.Handler
	server := &http.Server{Addr: serverUrl, Handler: router} ////TODO:Question какая разница между httpServer и http.ListenAndServe?

	go func() { //TODO:Question что за go func??
		log.Fatal(server.ListenAndServe())
	}()

	return server
}

//TODO:Question что за Chan??
func getKillSignalChan() chan os.Signal{
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}