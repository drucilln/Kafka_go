package cache

import (
	"sync"
	"untitled_folder/internal/model"
)

type Cache struct {
	Orders map[string]model.Order
	sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{Orders: make(map[string]model.Order)}
}

func (c *Cache) CacheSet(order model.Order) {
	c.Lock()
	defer c.Unlock()
	c.Orders[order.OrderUID] = order
}

func (c *Cache) CacheGet(orderUID string) (model.Order, bool) {
	c.RLock()
	defer c.RUnlock()
	order, ok := c.Orders[orderUID]
	return order, ok
}
