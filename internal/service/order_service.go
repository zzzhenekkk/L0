package service

import (
	"L0/internal/cache"
	"L0/internal/domain"
	"L0/internal/storage"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
)

type OrderService struct {
	repoDB *storage.PostgresOrderRepository
	cache  *cache.OrderCache
}

func NewOrderService(repo *storage.PostgresOrderRepository, cache *cache.OrderCache) *OrderService {
	os := &OrderService{repoDB: repo, cache: cache}
	err := os.LoadCache()
	if err != nil {
		log.Fatalf("Unable to load cache: %v\n", err)
	}
	return os
}

func (s *OrderService) ProcessOrder(msg *stan.Msg) {
	var order domain.Order

	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Printf("Error unmarshalling order: %v", err)
		return
	}

	err = s.repoDB.AddOrder(&order)
	if err != nil {
		log.Printf("Error adding order to database: %v", err)
		return
	}

	s.cache.Set(&order)

	if err := msg.Ack(); err != nil {
		log.Printf("Error acknowledging message: %v", err)
	} else {
		log.Println("Message acknowledged successfully")
	}
}

func (s *OrderService) LoadCache() error {
	err := s.repoDB.GetAllOrders(s.cache)
	if err != nil {
		return err
	}

	//s.cache.LoadFromDB(orders)
	return nil
}

func (s *OrderService) GetOrder(orderUID string) (*domain.Order, bool) {
	// Сначала пытаемся получить заказ из кэша
	if order, ok := s.cache.Get(orderUID); ok {
		return order, true
	}

	// Если в кэше нет, получаем из базы данных
	order, err := s.repoDB.GetOrder(orderUID)
	if err != nil {
		return nil, false
	}

	// Сохраняем в кэш
	s.cache.Set(order)
	return order, true
}
