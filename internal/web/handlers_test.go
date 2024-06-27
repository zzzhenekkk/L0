package web_test

import (
	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/domain"
	"L0/internal/service"
	"L0/internal/storage"
	"L0/internal/web"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func setupTestServer(t *testing.T) *web.Server {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
		},
		DB: config.DBConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "test_user",
			Password: "test_password",
			DBName:   "test_db",
			SSLMode:  "disable",
		},
	}
	repo := &storage.PostgresOrderRepository{}
	_, err := repo.ConnectDB(cfg)
	require.NoError(t, err)
	orderCache := cache.NewOrderCache()
	orderService := service.NewOrderService(repo, orderCache)

	// Создание таблиц, если они не существуют
	_, err = repo.DB.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS orders (
		order_uid VARCHAR(50) PRIMARY KEY,
		track_number VARCHAR(50),
		entry VARCHAR(10),
		delivery JSONB,
		payment JSONB,
		items JSONB,
		locale VARCHAR(10),
		internal_signature VARCHAR(50),
		customer_id VARCHAR(50),
		delivery_service VARCHAR(50),
		shardkey VARCHAR(10),
		sm_id INT,
		date_created TIMESTAMP,
		oof_shard VARCHAR(10)
	);
	
	CREATE TABLE IF NOT EXISTS uncorrect_orders_subscribe (
		id SERIAL PRIMARY KEY,
		order_data TEXT,
		sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`)
	require.NoError(t, err)

	return web.NewServer(cfg, orderService)
}

func TestGetOrderHandler(t *testing.T) {
	server := setupTestServer(t)

	order := &domain.Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: domain.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: domain.Payment{
			Transaction:  "b563feb7b2b84b6test",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []domain.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       time.Date(2021, 11, 26, 6, 22, 19, 0, time.UTC),
		OofShard:          "1",
	}

	orderData, err := json.Marshal(order)
	require.NoError(t, err)

	// Add order to the service
	err = server.GetOrderService().ProcessOrderData(orderData)
	require.NoError(t, err)

	req, err := http.NewRequest("GET", "/orders/b563feb7b2b84b6test", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	web.SetupRoutes(r, server.GetOrderService())
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var fetchedOrder domain.Order
	err = json.Unmarshal(rr.Body.Bytes(), &fetchedOrder)
	require.NoError(t, err)
	assert.Equal(t, order.OrderUID, fetchedOrder.OrderUID)
}

func TestGetOrderHandlerNotFound(t *testing.T) {
	server := setupTestServer(t)

	req, err := http.NewRequest("GET", "/orders/nonexistent", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	web.SetupRoutes(r, server.GetOrderService())
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
