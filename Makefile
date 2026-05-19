APP_NAME=customer-support-hub
MIGRATIONS_DIR=./migrations
DB_DSN=postgres://support:support123@localhost:5432/customer_support?sslmode=disable

.PHONY: setup dev build test lint tidy docker-up docker-down migrate-up migrate-down migrate-status dev-docker

setup:
	go mod tidy

dev: migrate-up
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

migrate-up:
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up; \
	else \
		docker run --rm -v "$(PWD)/migrations:/migrations" \
			migrate/migrate -path /migrations \
			-database "postgres://support:support123@host.docker.internal:5432/customer_support?sslmode=disable" up; \
	fi

migrate-down:
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down; \
	else \
		docker run --rm -v "$(PWD)/migrations:/migrations" \
			migrate/migrate -path /migrations \
			-database "postgres://support:support123@host.docker.internal:5432/customer_support?sslmode=disable" down; \
	fi

migrate-status:
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" version; \
	else \
		docker run --rm -v "$(PWD)/migrations:/migrations" \
			migrate/migrate -path /migrations \
			-database "postgres://support:support123@host.docker.internal:5432/customer_support?sslmode=disable" version; \
	fi

# Subir infra local e aplicar migrations antes de rodar a API

dev-docker: docker-up migrate-up dev
