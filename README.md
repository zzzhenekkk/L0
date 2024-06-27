# Сервис Обработки Заказов

Сервис обработки заказов - это демонстрационный сервис для управления и отображения данных о заказах. Сервис подключается к базе данных PostgreSQL, подписывается на сервер NATS Streaming для получения данных о заказах, кэширует данные в памяти и предоставляет HTTP API для получения информации о заказах по их идентификаторам.

## Возможности

- **Интеграция с PostgreSQL:** Хранение данных о заказах в базе данных PostgreSQL.
- **NATS Streaming:** Подписка на NATS Streaming для получения данных о заказах.
- **Кэширование в памяти:** Кэширование данных о заказах в памяти для быстрого доступа.
- **HTTP API:** Предоставление API для получения информации о заказах по их идентификаторам.
- **Устойчивость:** Восстановление кэша из базы данных в случае перезапуска сервиса.
- **Нагрузочное тестирование:** Поддержка нагрузочного тестирования с использованием WRK и Vegeta.

## Требования

- Go 1.19 или новее
- PostgreSQL
- NATS Streaming Server
- WRK (для нагрузочного тестирования)
- Vegeta (для нагрузочного тестирования)

## Конфигурация

1. Создайте файл конфигурации в директории `./configs` с именем `config.yaml`:

    ```yaml
    server:
      port: 8080

    db:
      host: localhost
      port: 5432
      user: your_db_user
      password: your_db_password
      dbname: your_db_name
      sslmode: disable

    nats:
      cluster: test-cluster
      client: test-client
      clientNotifier: test-client-notifier
      subject: orders
      durableName: my-durable
    ```

## Настройка базы данных

1. Создайте необходимые таблицы в вашей базе данных PostgreSQL:

    ```sql
    CREATE TABLE orders (
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

    CREATE TABLE uncorrect_orders_subscribe (
        id SERIAL PRIMARY KEY,
        order_data TEXT,
        sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
    ```

## Запуск сервиса

1. Запустите сервис обработки заказов:

    ```sh
    make service
    ```

2. Запустите сервис отправки уведомления:

    ```sh
    make notifier
    ```

## Нагрузочное тестирование

### Использование WRK

Для тестирования конечной точки публикации заказов с помощью WRK:

1. Создайте файл `post.lua` со следующим содержимым:

    ```lua
    wrk.method = "POST"
    wrk.body   = [[
    {
      "order_uid": "b563feb7b2b84b6test",
      "track_number": "WBILMTESTTRACK",
      "entry": "WBIL",
      "delivery": {
        "name": "Test Testov",
        "phone": "+9720000000",
        "zip": "2639809",
        "city": "Kiryat Mozkin",
        "address": "Ploshad Mira 15",
        "region": "Kraiot",
        "email": "test@gmail.com"
      },
      "payment": {
        "transaction": "b563feb7b2b84b6test",
        "currency": "USD",
        "provider": "wbpay",
        "amount": 1817,
        "payment_dt": 1637907727,
        "bank": "alpha",
        "delivery_cost": 1500,
        "goods_total": 317,
        "custom_fee": 0
      },
      "items": [
        {
          "chrt_id": 9934930,
          "track_number": "WBILMTESTTRACK",
          "price": 453,
          "rid": "ab4219087a764ae0btest",
          "name": "Mascaras",
          "sale": 30,
          "size": "0",
          "total_price": 317,
          "nm_id": 2389212,
          "brand": "Vivienne Sabo",
          "status": 202
        }
      ],
      "locale": "en",
      "internal_signature": "",
      "customer_id": "test",
      "delivery_service": "meest",
      "shardkey": "9",
      "sm_id": 99,
      "date_created": "2021-11-26T06:22:19Z",
      "oof_shard": "1"
    }
    ]]
    wrk.headers["Content-Type"] = "application/json"
    ```

2. Запустите нагрузочный тест:

    ```sh
    make testWRK
    ```

### Использование Vegeta

Для тестирования конечной точки публикации заказов с помощью Vegeta:

1. Запустите нагрузочный тест:

    ```sh
    make testVegeta
    ```

Это создаст отчет и график результатов.

## Запуск тестов

Чтобы запустить модульные тесты:

```sh
make test
