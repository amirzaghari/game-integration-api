APP_NAME=app

.PHONY: build run test migrate-up migrate-down up down

up:
	docker-compose up -d

down:
	docker-compose down

build:
	docker-compose exec app go build -o $(APP_NAME) ./cmd

run:
	docker-compose exec app go run ./cmd

test:
	docker-compose exec app go test ./...

migrate-up:
	docker-compose exec app migrate -path ./migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" up

migrate-down:
	docker-compose exec app migrate -path ./migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" down 