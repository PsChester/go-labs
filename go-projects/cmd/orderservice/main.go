package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"go-projects/pkg/orderservice/transport"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//Транзакции
//tx, error := db.Begin()
//если ошибка -> log.Fatal(error)
//defer tx.
//tx.Query()
//tx.Commit()
//Нужно чистить память, напр: rows.Close() У структур, которые создаёшь сам, чистить не нужно
//А вот сырые данные, которые приходят откуда-то, напр. из database или body в post-запросе, нужно подчистить

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("orderservice.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666) //TODO:Question Как работает?
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	//Подключение базы данных
	database, err := sql.Open("mysql", `root:root@/orderservice`)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatal(err)
	}

	serverUrl := ":8000"
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")
	killSignalChan := getKillSignalChannel()
	server := startServer(serverUrl)

	waitForKillSignal(killSignalChan)
	server.Shutdown(context.Background())
}

func startServer(serverUrl string) *http.Server {
	router := transport.Router()                             //http.Handler
	server := &http.Server{Addr: serverUrl, Handler: router} ////TODO:Question какая разница между httpServer и http.ListenAndServe?

	go func() { //TODO:Question что за go func??
		log.Fatal(server.ListenAndServe())
	}()

	return server
}

//TODO:Question что за Chan??
func getKillSignalChannel() chan os.Signal {
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
