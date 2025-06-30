APP_NAME=app

.PHONY: build run test up down

up:
	docker-compose up -d

down:
	docker-compose down

build:
	cp .env.example .env
	docker-compose up --build

local-dev:
	export $$(grep -v '^#' .env | xargs) && air

run:
	docker-compose exec app go run ./cmd

test:
	docker-compose exec app go test ./test/...
