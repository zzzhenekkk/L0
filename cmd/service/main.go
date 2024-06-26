package main

import (
	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/nats"
	"L0/internal/service"
	"L0/internal/storage"
	"L0/internal/web"
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

	// подключение к nats streaming
	//sc, err := stan.Connect(cfg.Nats.Cluster, cfg.Nats.Client)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer sc.Close()

	// подключение к БД
	repo := &storage.PostgresOrderRepository{}
	_, err = repo.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer repo.CloseDB()

	//subscribeOrder(sc, cfg, repo)
	orderCache := cache.NewOrderCache()
	defer orderCache.PrintCache()
	orderService := service.NewOrderService(repo, orderCache)

	natsSubscriber := nats.NewNatsSubscriber(cfg, orderService)
	err = natsSubscriber.SubscribeToOrders()
	if err != nil {
		log.Fatalf("Unable to subscribe to orders: %v\n", err)
	}

	server := web.NewServer(cfg, orderService)
	server.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

}

//
//func subscribeOrder(sc stan.Conn, cfg *config.Config, repo *storage.PostgresOrderRepository) {
//
//	subscriptionOptions := []stan.SubscriptionOption{
//		stan.DurableName(cfg.Nats.DurableName),
//		stan.SetManualAckMode(),
//	}
//
//	_, err := sc.Subscribe(cfg.Nats.Subject, processing, subscriptionOptions...)
//
//	if err != nil {
//		log.Println("Error subscribing to orders: %v", err)
//		return
//	}
//	log.Println("Subscribed to orders")
//}
//
//func processing(msg *stan.Msg) {
//	var order domain.Order
//	err := json.Unmarshal(msg.Data, &order)
//	if err != nil {
//		log.Println("Error unmarshalling order: %v", err)
//		return
//	}
//	log.Println("Received order: %+v", order)
//
//	if err := msg.Ack(); err != nil {
//		log.Printf("Error acknowledging message: %v", err)
//	} else {
//		log.Println("Message acknowledged successfully")
//	}
//}
