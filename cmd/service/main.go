package main

import (
	"L0/internal/config"
	"L0/internal/domain"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"os/signal"
)

func main() {
	//if pathEx, err := os.Getwd(); err == nil {
	//	println(pathEx)
	//}

	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Unmarshal config %+v", conf)

	sc, err := stan.Connect(cfg.Nats.Cluster, cfg.Nats.Client)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	subscribeOrder(sc, cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

}

func subscribeOrder(sc stan.Conn, cfg *config.Config) {

	subscriptionOptions := []stan.SubscriptionOption{
		stan.DurableName(cfg.Nats.DurableName),
		stan.SetManualAckMode(),
	}

	_, err := sc.Subscribe(cfg.Nats.Subject, processing, subscriptionOptions...)

	if err != nil {
		log.Println("Error subscribing to orders: %v", err)
		return
	}
	log.Println("Subscribed to orders")
}

func processing(msg *stan.Msg) {
	var order domain.Order
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Println("Error unmarshalling order: %v", err)
		return
	}
	log.Println("Received order: %+v", order)

	if err := msg.Ack(); err != nil {
		log.Printf("Error acknowledging message: %v", err)
	} else {
		log.Println("Message acknowledged successfully")
	}
}
