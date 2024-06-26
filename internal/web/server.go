package web

import (
	"L0/internal/config"
	"L0/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, orderService *service.OrderService) *Server {
	r := mux.NewRouter()
	SetupRoutes(r, orderService)
	adr := ":" + strconv.Itoa(cfg.Server.Port)
	return &Server{
		httpServer: &http.Server{
			Handler:      r,
			Addr:         adr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
	}
}

func (s *Server) Run() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", s.httpServer.Addr, err)
		}
	}()
	log.Printf("Server is running on port %s", s.httpServer.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	if err := s.httpServer.Close(); err != nil {
		log.Fatalf("Server Close: %v", err)
	}
}
