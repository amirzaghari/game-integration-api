APP_NAME=app

.PHONY: help up down build run test local-dev

help:
	@echo ""
	@echo "▄▖         ▄▖  ▗         ▗ ▘      ▄▖▄▖▄▖"
	@echo "▌ ▀▌▛▛▌█▌  ▐ ▛▌▜▘█▌▛▌▛▘▀▌▜▘▌▛▌▛▌  ▌▌▙▌▐ "
	@echo "▙▌█▌▌▌▌▙▖  ▟▖▌▌▐▖▙▖▙▌▌ █▌▐▖▌▙▌▌▌  ▛▌▌ ▟▖"
	@echo "                   ▄▌                   "
	@echo ""
	@echo ""
	@echo "Makefile Commands:"
	@echo ""
	@echo "  make up         Start all services (db, wallet, app) via Docker Compose"
	@echo "  make build      Build the app container and copy .env.example to .env if present"
	@echo "  make run        Run the app inside the container"
	@echo "  make down       Stop all services"
	@echo "  make test       Run all tests in the test/ directory inside the container"
	@echo "  make local-dev  Run the app locally with hot reload (requires air)"
	@echo ""

.DEFAULT_GOAL := help

up:
	docker-compose up -d

down:
	docker-compose down

build:
	cp .env.example .env || true
	docker-compose up --build

local-dev:
	export $$(grep -v '^#' .env | xargs) && air

run:
	docker-compose exec app go run ./cmd

test:
	docker-compose exec app go test ./test/...
