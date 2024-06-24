package main

import (
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Connect to NATS Streaming Server
	sc, err := stan.Connect("test-cluster", "client-123")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// Simple Publisher
	go func() {
		for {
			sc.Publish("my-channel", []byte("Hello NATS Streaming!"))
			time.Sleep(5 * time.Second)
		}
	}()

	// Simple Subscriber
	sub, err := sc.Subscribe("my-channel", func(m *stan.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// Run forever
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit
}
