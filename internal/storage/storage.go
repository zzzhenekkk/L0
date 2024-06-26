package storage

import (
	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type PostgresOrderRepository struct {
	DB *pgx.Conn
}

func (db *PostgresOrderRepository) ConnectDB(cfg *config.Config) (*PostgresOrderRepository, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	db.DB = conn
	return db, nil
}

func (repo *PostgresOrderRepository) CloseDB() {
	if repo.DB != nil {
		repo.DB.Close(context.Background())
	}
}

func (repo *PostgresOrderRepository) AddOrder(order *domain.Order) error {
	query := `INSERT INTO orders (
        order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
    ) ON CONFLICT (order_uid) DO UPDATE
    SET
        track_number = EXCLUDED.track_number,
        entry = EXCLUDED.entry,
        delivery = EXCLUDED.delivery,
        payment = EXCLUDED.payment,
        items = EXCLUDED.items,
        locale = EXCLUDED.locale,
        internal_signature = EXCLUDED.internal_signature,
        customer_id = EXCLUDED.customer_id,
        delivery_service = EXCLUDED.delivery_service,
        shardkey = EXCLUDED.shardkey,
        sm_id = EXCLUDED.sm_id,
        date_created = EXCLUDED.date_created,
        oof_shard = EXCLUDED.oof_shard;`

	deliveryJSON, err := json.Marshal(order.Delivery)
	if err != nil {
		return err
	}

	paymentJSON, err := json.Marshal(order.Payment)
	if err != nil {
		return err
	}

	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return err
	}

	_, err = repo.DB.Exec(context.Background(), query,
		order.OrderUID, order.TrackNumber, order.Entry, deliveryJSON, paymentJSON, itemsJSON, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	return err
}

func (repo *PostgresOrderRepository) GetOrder(orderUID string) (*domain.Order, error) {
	query := `SELECT order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
              FROM orders WHERE order_uid = $1`

	var order domain.Order
	var deliveryJSON, paymentJSON, itemsJSON []byte

	err := repo.DB.QueryRow(context.Background(), query, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &deliveryJSON, &paymentJSON, &itemsJSON,
		&order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.ShardKey,
		&order.SmID, &order.DateCreated, &order.OofShard)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(deliveryJSON, &order.Delivery)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(paymentJSON, &order.Payment)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(itemsJSON, &order.Items)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (repo *PostgresOrderRepository) GetAllOrders(cache *cache.OrderCache) error {
	query := `SELECT order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders`

	rows, err := repo.DB.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var order domain.Order
		var deliveryJSON, paymentJSON, itemsJSON []byte

		err = rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &deliveryJSON, &paymentJSON, &itemsJSON,
			&order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.ShardKey,
			&order.SmID, &order.DateCreated, &order.OofShard)

		if err != nil {
			return err
		}

		err = json.Unmarshal(deliveryJSON, &order.Delivery)
		if err != nil {
			return err
		}

		err = json.Unmarshal(paymentJSON, &order.Payment)
		if err != nil {
			return err
		}

		err = json.Unmarshal(itemsJSON, &order.Items)
		if err != nil {
			return err
		}

		cache.Set(&order)
	}

	return nil
}

func (repo *PostgresOrderRepository) SaveUncorrectOrder(order string) error {
	query := `INSERT INTO uncorrect_orders_subscribe (order_data) VALUES ($1)`
	_, err := repo.DB.Exec(context.Background(), query, order)
	return err
}
