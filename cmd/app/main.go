package main

import (
	"Kafka_go/internal/cache"
	"Kafka_go/internal/handler"
	"Kafka_go/internal/kafka"
	"Kafka_go/internal/postgres"
	"fmt"
	"log"
	"net/http"
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
	go kafka.InitKafka(c, db)
	http.HandleFunc("/order", handler.GetOrderHandler(c))

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
