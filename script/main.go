package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
)

func main() {
	w := kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: "orders",
	}

	defer w.Close()

	file := os.Args[1]

	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	err = w.WriteMessages(context.Background(), kafka.Message{
		Value: data,
	})
	if err != nil {
		log.Fatal(err)
	}

}
