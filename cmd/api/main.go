package main

import (
	"fmt"
	"log"
	"math/rand"
	"smarthome-home/internal"
	"time"
)

func main() {
	log.Println("STARTING smarthome")
	log.SetFlags(log.LstdFlags | log.Llongfile)
	var err error
	config := initConfig()
	server := initServer(config.Server.Name, config.Server.MaxRequestBodySize)

	var storage *Storage
	sa := &internal.ServicesAvailable{}
	if storage, err = initStorage(config.DBMySQL, "mysql", sa); err != nil {
		log.Fatal(err)
	}
	var push *Push
	if push, err = initPusher(config.PushMQTT, "paho", sa); err != nil {
		log.Fatal(err)
	}
	var service *Service
	if service, err = initService(storage, push, config.Security.Hash); err != nil {
		log.Fatal(err)
	}
	var handler *Handler
	if handler, err = initHandler(service, sa); err != nil {
		log.Fatal(err)
	}

	go func(sa *internal.ServicesAvailable) {
		for !sa.DB {
			time.Sleep(1 * time.Second)
			log.Println("TRYING AGAIN")
		}
		if err := service.user.Create("Chiel Timmermans", "chieltimmermans@hotmail.com", "test1234", "admin"); err != nil {
			log.Println(err)
		}
	}(sa)

	rand.Seed(time.Now().UnixNano())

	initRouter(server, handler, config.Router, config.CORS)

	log.Printf("Starting up smarthome-home back-end, listening on port: %d\n", config.Server.Port)
	if err := server.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
