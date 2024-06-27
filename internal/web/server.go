package web

import (
	"L0/internal/config"
	"L0/internal/service"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	httpServer   *http.Server
	orderService *service.OrderService
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
		orderService: orderService,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) GetOrderService() *service.OrderService {
	return s.orderService
}
