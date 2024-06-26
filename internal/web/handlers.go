package web

import (
	"L0/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router, orderService *service.OrderService) {
	r.HandleFunc("/orders/{id}", getOrderHandler(orderService)).Methods("GET")
}

func getOrderHandler(orderService *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderUID := vars["id"]

		order, found := orderService.GetOrder(orderUID)
		if !found {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	}
}
