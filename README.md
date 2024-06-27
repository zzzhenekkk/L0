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
Для тестирования конечной точки публикации заказов с помощью WRK

1. Запустите нагрузочный тест:

    ```sh
    make testWRK
    ```

### Использование Vegeta

Для тестирования конечной точки публикации заказов с помощью Vegeta:

1. Запустите нагрузочный тест:

    ```sh
    make testVegeta
    ```

Это создаст отчет.

## Запуск тестов

Чтобы запустить модульные тесты:

```sh
make test
