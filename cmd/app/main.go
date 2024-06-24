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
	conf, err := config.LoadConfig("../../configs")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Unmarshal config %+v", conf)

	sc, err := stan.Connect(conf.Nats.Cluster, "order-service-client")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	subscribeOrder(sc)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

}

func subscribeOrder(sc stan.Conn) {
	_, err := sc.Subscribe("orders", func(msg *stan.Msg) {
		var order domain.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println("Error unmarshalling order: %v", err)
			return
		}
		log.Println("Received order: %+v", order)
	}, stan.DurableName("order-service-subscription"))

	if err != nil {
		log.Println("Error subscribing to orders: %v", err)
		return
	}
	log.Println("Subscribed to orders")
}
