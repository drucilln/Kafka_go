package main

import (
	"fmt"
	"log"
	"net/http"
	"untitled_folder/internal/cache"
	"untitled_folder/internal/handler"
	"untitled_folder/internal/nats"
	"untitled_folder/internal/postgres"
)

func main() {
	db, err := postgres.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	c := cache.NewCache()
	if err := postgres.LoadCacheFromDB(db, c); err != nil {
		log.Fatalf("Error loading cache from DB: %v", err)
	}
	//fmt.Println(c.Orders)
	go nats.InitKafka(c, db)
	http.HandleFunc("/order", handler.GetOrderHandler(c))

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
