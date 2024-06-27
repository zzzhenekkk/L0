# Define the packages to be tested
TEST_PACKAGES := ./internal/storage ./internal/service ./internal/web

all: service

test:
	@echo "Running tests..."
	@go test $(TEST_PACKAGES) -v

service:
	go run cmd/service/main.go

notifier:
	go run cmd/notifier/main.go

testWRK:
	wrk -t12 -c400 -d10s http://localhost:8080/orders/b563feb7b2b84b6test
	wrk -t12 -c400 -d10s -s post.lua http://localhost:8081/publish

notifierPost:
	go run cmd/notifierPost/main.go


testVegeta:
	echo "GET http://localhost:8080/orders/b563feb7b2b84b6test" | vegeta attack -duration=10s -rate=100 | tee results.bin | vegeta report

