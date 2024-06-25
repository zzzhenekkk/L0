package main

import (
	"L0/internal/config"
	"L0/internal/domain"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"os"
)

func main() {
	conf, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Unmarshal config %+v", conf)

	sc, err := stan.Connect(conf.Nats.Cluster, conf.Nats.ClientNotifier)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	publicOneOrder(sc, conf)

}

func publishOrder(sc stan.Conn, order domain.Order) {
	data, err := json.Marshal(order)
	if err != nil {
		log.Println("Error marshalling order: %v", err)
	}

	err = sc.Publish("orders", data)
	if err != nil {
		log.Println("Error publishing order: %v", err)
		return
	}
	log.Println("Order published successfully", order.OrderUID)
}

func publicOneOrder(sc stan.Conn, cfg *config.Config) {
	data, err := os.ReadFile("./internal/notifier/correct.json")

	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	log.Printf("Expected json: %s", data)

	err = sc.Publish(cfg.Nats.Subject, data)
	if err != nil {
		log.Println("Error publishing order: %v", err)
		return
	}
	log.Println("Order published successfully")
}
