package main

import (
	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/nats"
	"L0/internal/service"
	"L0/internal/storage"
	"L0/internal/web"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal(err)
	}

	repo := &storage.PostgresOrderRepository{}
	_, err = repo.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer repo.CloseDB()

	orderCache := cache.NewOrderCache()
	defer orderCache.PrintCache()
	orderService := service.NewOrderService(repo, orderCache)

	natsSubscriber := nats.NewNatsSubscriber(cfg, orderService)
	err = natsSubscriber.SubscribeToOrders()
	if err != nil {
		log.Fatalf("Unable to subscribe to orders: %v\n", err)
	}
	defer natsSubscriber.Close()

	server := web.NewServer(cfg, orderService)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr(), err)
		}
	}()
	log.Printf("Server is running on port %s", server.Addr())

	<-ctx.Done()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Application stopped")
}
