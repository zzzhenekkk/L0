package main

import (
	"L0/internal/config"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
)

var (
	sc   stan.Conn
	conf *config.Config
)

func main() {
	var err error
	conf, err = config.LoadConfig("./configs")
	if err != nil {
		log.Fatal(err)
	}

	sc, err = stan.Connect(conf.Nats.Cluster, conf.Nats.ClientNotifier)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	http.HandleFunc("/publish", publishOrderHandler)

	log.Println("Notifier server is running on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

func publishOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Failed to marshal order", http.StatusInternalServerError)
		return
	}

	err = sc.Publish(conf.Nats.Subject, data)
	if err != nil {
		http.Error(w, "Failed to publish order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order published successfully"))
}
