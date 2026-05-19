APP_NAME=customer-support-hub

.PHONY: setup dev build test lint tidy docker-up docker-down

setup:
	go mod tidy

dev:
	go run ./cmd/api

build:
	go build -o bin/$(APP_NAME) ./cmd/api

test:
	go test ./...

lint:
	gofmt -w ./cmd ./internal
	go vet ./...

tidy:
	go mod tidy

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down
