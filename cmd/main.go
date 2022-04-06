package main

import (
	"github.com/go-playground/validator/v10"
	"log"

	"L0/internal/services"
)

var db *services.DB
var cache *services.Cache
var service *services.Service

func init() {
	var err error
	connString := "user=intern password=qwerty1234 host=localhost port=5432 dbname=wb sslmode=disable"
	db, err = services.NewDB(connString)
	if err != nil {
		log.Fatal(err)
	}

	cache = services.NewCache()
	err = cache.DataRecovery(db)
	if err != nil {
		log.Fatal(err)
	}

	valid := validator.New()

	service = services.NewService(cache, db, valid)
	err = service.ConnectionToNats()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer db.Close()
	defer service.Client.Close()

	err := service.Subscribe()
	if err != nil {
		log.Fatal(err)
	}
	defer service.Sub.Unsubscribe()

	service.Run()
}
