package main

//TODO:Question как настроить создание исполняемых файлов в bin?

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderservice/pkg/orderservice/model"
	"orderservice/pkg/orderservice/transport"
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
//А вот сырые данные, которые приходят откуда-то, напр. из database.mysql или body в post-запросе, нужно подчистить

//TODO:Question как хранить миграции базы данных?
//TODO:Question как обращаться к базе данных не объявля методы для типа OrderService
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("../../orderservice.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666) //TODO:Question os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666?
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	//Подключение базы данных
	database, err := sql.Open("mysql", `root:root@/orderservice?parseTime=true`)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatal(err)
	}

	//TODO:мб объеденить с server и serverUrl
	databaseServer := model.OrderService{Database: database}

	serverUrl := ":8000"
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")
	killSignalChan := getKillSignalChannel()
	server := startServer(serverUrl, databaseServer)

	waitForKillSignal(killSignalChan)
	server.Shutdown(context.Background())
}

func startServer(serverUrl string, databaseServer model.OrderService) *http.Server {
	router := transport.Router(&databaseServer)              //http.Handler
	server := &http.Server{Addr: serverUrl, Handler: router} ////TODO:Question какая разница между httpServer и http.ListenAndServe?

	go func() { //Answer go func - это goroutine, функция выполняющая ассинхронно
		log.Fatal(server.ListenAndServe())
	}()

	return server
}

//Answer: сhan - тип называемый "канал" (channel). По сути обычная очередь //TODO:Question или я ошибаюсь?
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
