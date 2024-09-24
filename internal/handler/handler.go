package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"untitled_folder/internal/cache"
)

func GetOrderHandler(c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderId := r.URL.Query().Get("id")
		if orderId == "" {
			log.Println("Missing id parameter")
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}
		order, exists := c.CacheGet(orderId)
		if !exists {
			log.Printf("Order with id %s not found", orderId)
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(order); err != nil {
			log.Printf("Error encoding json: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
