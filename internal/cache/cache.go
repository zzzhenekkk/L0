package cache

import (
	"L0/internal/domain"
	"sync"
)

type OrderCache struct {
	cache sync.Map
}

func NewOrderCache() *OrderCache {
	return &OrderCache{}
}

func (c *OrderCache) Get(orderUID string) (*domain.Order, bool) {
	order, ok := c.cache.Load(orderUID)
	if !ok {
		return nil, false
	}
	return order.(*domain.Order), true
}

func (c *OrderCache) Set(order *domain.Order) {
	c.cache.Store(order.OrderUID, order)
}
