#!make
include .env
DB_URL=postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable


.PHONY: all build run unittest clean migrate-up migrate-down


# Build the application
all: build

build:
	@echo "Building..."
	@sqlc generate
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run main.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
unittest:
	@echo "Testing..."
	@go test -v  $$(go list ./... | grep -v e2e)

e2e:
	@echo "E2E Testing..."
	@go test -v ./tests/e2e/...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main


migrate-up:
	@echo "Up migrating..."
	@echo "$(DB_URL)"
	@goose -dir sql/schema postgres "$(DB_URL)" up

migrate-down:
	@echo "Down migrating..."
	@goose -dir sql/schema postgres "$(DB_URL)" down

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

