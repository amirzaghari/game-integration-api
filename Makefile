APP_NAME=game-integration-api

.PHONY: build run test migrate-up migrate-down docker-up docker-down

build:
	go build -o $(APP_NAME) ./cmd

run:
	go run ./cmd

test:
	go test ./...

migrate-up:
	migrate -path ./migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" down

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down 