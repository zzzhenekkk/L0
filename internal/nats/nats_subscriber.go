package nats

import (
	"L0/internal/config"
	"L0/internal/service"
	"github.com/nats-io/stan.go"
	"log"
)

type NatsSubscriber struct {
	cfg          *config.Config
	orderService *service.OrderService
	natsConn     stan.Conn
}

func NewNatsSubscriber(cfg *config.Config, orderService *service.OrderService) *NatsSubscriber {
	return &NatsSubscriber{cfg: cfg, orderService: orderService}
}

func (n *NatsSubscriber) SubscribeToOrders() error {
	sc, err := stan.Connect(n.cfg.Nats.Cluster, n.cfg.Nats.Client)
	if err != nil {
		return err
	}
	n.natsConn = sc

	subscriptionOptions := []stan.SubscriptionOption{
		stan.DurableName(n.cfg.Nats.DurableName),
		stan.SetManualAckMode(),
	}

	_, err = sc.Subscribe(n.cfg.Nats.Subject, n.orderService.ProcessOrder, subscriptionOptions...)
	if err != nil {
		return err
	}

	log.Println("Subscribed to orders")
	return nil
}

func (n *NatsSubscriber) Close() {
	if n.natsConn != nil {
		n.natsConn.Close()
	}
}
