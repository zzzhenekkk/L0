package storage_test

import (
	"L0/internal/config"
	"L0/internal/domain"
	"L0/internal/storage"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *storage.PostgresOrderRepository {
	cfg := &config.Config{
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

	// Создание таблиц, если они не существуют
	_, err = repo.DB.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS orders (
		order_uid VARCHAR(100) PRIMARY KEY,
		track_number VARCHAR(100),
		entry VARCHAR(100),
		delivery JSONB,
		payment JSONB,
		items JSONB,
		locale VARCHAR(100),
		internal_signature VARCHAR(100),
		customer_id VARCHAR(100),
		delivery_service VARCHAR(100),
		shardkey VARCHAR(100),
		sm_id INT,
		date_created TIMESTAMP,
		oof_shard VARCHAR(100)
	);
	
	CREATE TABLE IF NOT EXISTS uncorrect_orders_subscribe (
		id SERIAL PRIMARY KEY,
		order_data TEXT,
		sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`)
	require.NoError(t, err)

	return repo
}

func TestAddOrder(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.CloseDB()

	order := &domain.Order{
		OrderUID:    "1234567890",
		TrackNumber: "TRACK12345",
		Entry:       "WBIL",
		Delivery: domain.Delivery{
			Name:    "Test Name",
			Phone:   "+1234567890",
			Zip:     "123456",
			City:    "Test City",
			Address: "Test Address",
			Region:  "Test Region",
			Email:   "test@example.com",
		},
		Payment: domain.Payment{
			Transaction:  "test_transaction",
			Currency:     "USD",
			Provider:     "test_provider",
			Amount:       1000,
			PaymentDT:    time.Now().Unix(),
			Bank:         "test_bank",
			DeliveryCost: 100,
			GoodsTotal:   900,
			CustomFee:    0,
		},
		Items: []domain.Item{
			{
				ChrtID:      123,
				TrackNumber: "test_track",
				Price:       500,
				Rid:         "test_rid",
				Name:        "Test Item",
				Sale:        50,
				Size:        "M",
				TotalPrice:  450,
				NmID:        12345,
				Brand:       "Test Brand",
				Status:      1,
			},
		},
		Locale:            "en",
		InternalSignature: "signature",
		CustomerID:        "customer_123",
		DeliveryService:   "delivery_svc",
		ShardKey:          "shard",
		SmID:              1,
		DateCreated:       time.Now(),
		OofShard:          "oof", // Ensure this fits within character varying(10)
	}

	err := repo.AddOrder(order)
	assert.NoError(t, err)

	fetchedOrder, err := repo.GetOrder("1234567890")
	assert.NoError(t, err)
	assert.Equal(t, order.OrderUID, fetchedOrder.OrderUID)
	assert.Equal(t, order.TrackNumber, fetchedOrder.TrackNumber)
}

func TestSaveUncorrectOrder(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.CloseDB()

	err := repo.SaveUncorrectOrder("incorrect order data")
	assert.NoError(t, err)
}
