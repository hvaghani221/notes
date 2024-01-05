#!make
include .env
DB_URL=postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable


.PHONY: all build run unittest e2e clean tools migrate-up migrate-down


# Build the application
all: build

build:
	@echo "Building..."
	@sqlc generate
	@go build -o main main.go

# Run the application
run:
	@go run main.go

# Create DB container
db-up:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up -d; \
	fi

# Shutdown DB container
db-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
unittest:
	@echo "Testing..."
	@go clean -testcache && go test -v  $$(go list ./... | grep -v e2e)

e2e:
	@chmod +x e2e.sh
	@./e2e.sh

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

tools:
	@echo "Installing tools..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest

sqlc:
	@echo "Generating SQL queries..."
	@sqlc generate

migrate-up:
	@echo "Up migrating..."
	@goose -dir sql/schema postgres "$(DB_URL)" up

migrate-down:
	@echo "Down migrating..."
	@goose -dir sql/schema postgres "$(DB_URL)" down
