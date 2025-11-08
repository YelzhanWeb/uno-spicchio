.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $1, $2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/restaurant-crm cmd/app

run: ## Run the application
	go run cmd/app/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users)
	migrate create -ext sql -dir migrations -seq $(name)

deps: ## Download dependencies
	go mod download
	go mod tidy

lint: ## Run linter
	golangci-lint run

.DEFAULT_GOAL := help