package kafka

import (
	"Kafka_go/internal/cache"
	"Kafka_go/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/xeipuuv/gojsonschema"
	"gorm.io/gorm"
	"log"
	"path/filepath"
)

func InitKafka(c *cache.Cache, db *gorm.DB) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "orders_server",
		Topic:   "orders",
	})

	defer r.Close()

	handleMsg(r, c, db)
}

func handleMsg(r *kafka.Reader, c *cache.Cache, db *gorm.DB) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	absPath, err := filepath.Abs("../../internal/model/schema.json")
	if err != nil {
		log.Fatal(err)
	}

	schema := gojsonschema.NewReferenceLoader("file://" + absPath)

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				break
			}
			log.Printf("Error reading message: %s\n", err)
			continue
		}

		dataLoader := gojsonschema.NewBytesLoader(m.Value)

		result, err := gojsonschema.Validate(schema, dataLoader)
		if err != nil {
			log.Printf("Error when validating JSON: %s\n", err)
			continue
		}

		if !result.Valid() {
			log.Printf("Validation errors:")
			for _, desc := range result.Errors() {
				log.Printf("- %s\n", desc)
			}
			continue
		}

		var order model.Order
		err = json.Unmarshal(m.Value, &order)
		if err != nil {
			log.Printf("Error unmarshalling message: %s\n", err)
			continue
		}

		tx := db.Begin()
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			if pgErr, ok := err.(*pq.Error); ok {
				if pgErr.Code == "23505" {
					log.Printf("Order already exists: %s\n", order)
				} else {
					log.Printf("Error creating order: %s\n", err)
				}
			}
			log.Printf("Error creating order: %v", err)
			continue
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			log.Printf("Error committing order: %v", err)
			continue
		}

		c.CacheSet(order)
		fmt.Printf("Order: %+v\n", order)
	}
}
