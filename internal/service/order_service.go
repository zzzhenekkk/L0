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
		s.saveUncorrectOrder(msg.Data)
		msg.Ack()
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
		return
	}

	log.Println("Adding an order to the database successfully.", "OrderUID:", order.OrderUID)
}

func (s *OrderService) LoadCache() error {
	err := s.repoDB.GetAllOrders(s.cache)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetOrder(orderUID string) (*domain.Order, bool) {
	if order, ok := s.cache.Get(orderUID); ok {
		return order, true
	}

	order, err := s.repoDB.GetOrder(orderUID)
	if err != nil {
		return nil, false
	}

	s.cache.Set(order)
	return order, true
}

func (s *OrderService) saveUncorrectOrder(orderData []byte) {
	orderStr := string(orderData)
	if err := s.repoDB.SaveUncorrectOrder(orderStr); err != nil {
		log.Printf("Error saving uncorrect order: %v", err)
	} else {
		log.Println("Uncorrect order saved to database:", string(orderData))
	}
}
