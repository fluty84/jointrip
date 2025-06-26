# JoinTrip Makefile

# Variables
APP_NAME=jointrip
DOCKER_COMPOSE=docker-compose
GO_CMD=go
MIGRATE_CMD=migrate
DB_URL=postgres://postgres:jointrip_password@localhost:5432/jointrip?sslmode=disable

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help setup db-up db-down db-reset migrate-up migrate-down migrate-create run test test-verbose clean build build-frontend build-quick dev dev-full docker-build docker-run

# Default target
help: ## Show this help message
	@echo "$(BLUE)JoinTrip Development Commands$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Setup development environment
setup: ## Setup development environment
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	@cp .env.example .env || echo "$(YELLOW).env.example not found, using existing .env$(NC)"
	@$(GO_CMD) mod download
	@$(GO_CMD) mod tidy
	@echo "$(GREEN)Development environment setup complete!$(NC)"

# Database commands
db-up: ## Start database containers
	@echo "$(YELLOW)Starting database containers...$(NC)"
	@$(DOCKER_COMPOSE) up -d postgres redis
	@echo "$(GREEN)Database containers started!$(NC)"
	@echo "$(BLUE)Waiting for database to be ready...$(NC)"
	@sleep 5
	@$(DOCKER_COMPOSE) exec postgres pg_isready -U postgres -d jointrip || echo "$(YELLOW)Database might still be starting up$(NC)"

db-down: ## Stop database containers
	@echo "$(YELLOW)Stopping database containers...$(NC)"
	@$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Database containers stopped!$(NC)"

db-reset: ## Reset database (stop, remove volumes, start)
	@echo "$(YELLOW)Resetting database...$(NC)"
	@$(DOCKER_COMPOSE) down -v
	@$(DOCKER_COMPOSE) up -d postgres redis
	@echo "$(GREEN)Database reset complete!$(NC)"

db-logs: ## Show database logs
	@$(DOCKER_COMPOSE) logs -f postgres

# Migration commands
migrate-up: ## Run database migrations up
	@echo "$(YELLOW)Running migrations up...$(NC)"
	@$(MIGRATE_CMD) -path migrations -database "$(DB_URL)" up
	@echo "$(GREEN)Migrations completed!$(NC)"

migrate-down: ## Run database migrations down
	@echo "$(YELLOW)Running migrations down...$(NC)"
	@$(MIGRATE_CMD) -path migrations -database "$(DB_URL)" down
	@echo "$(GREEN)Migrations rolled back!$(NC)"

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "$(RED)Error: Please provide a migration name. Usage: make migrate-create name=migration_name$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Creating migration: $(name)$(NC)"
	@$(MIGRATE_CMD) create -ext sql -dir migrations $(name)
	@echo "$(GREEN)Migration created!$(NC)"

# Application commands
run: ## Run the application
	@echo "$(YELLOW)Starting JoinTrip API...$(NC)"
	@$(GO_CMD) run main.go embed.go

dev: db-up migrate-up run ## Start development environment (database + migrations + app)

dev-full: ## Start full development environment with frontend build
	@./scripts/dev.sh

build-frontend: ## Build React frontend
	@echo "$(YELLOW)Building React frontend...$(NC)"
	@cd web && source ~/.nvm/nvm.sh && npm run build
	@echo "$(GREEN)Frontend build complete!$(NC)"

build: build-frontend ## Build the full application (frontend + backend)
	@echo "$(YELLOW)Building Go backend with embedded frontend...$(NC)"
	@$(GO_CMD) build -o bin/$(APP_NAME) main.go embed.go
	@echo "$(GREEN)Build complete! Binary: bin/$(APP_NAME)$(NC)"

build-quick: ## Quick build using build script
	@./scripts/build.sh

# Testing commands
test: ## Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	@$(GO_CMD) test ./...
	@echo "$(GREEN)Tests completed!$(NC)"

test-verbose: ## Run tests with verbose output
	@echo "$(YELLOW)Running tests with verbose output...$(NC)"
	@$(GO_CMD) test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	@$(GO_CMD) test -coverprofile=coverage.out ./...
	@$(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Docker commands
docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(NC)"
	@docker build -t $(APP_NAME):latest .
	@echo "$(GREEN)Docker image built!$(NC)"

docker-run: ## Run application in Docker
	@echo "$(YELLOW)Running application in Docker...$(NC)"
	@docker run --rm -p 8080:8080 --env-file .env $(APP_NAME):latest

# Utility commands
clean: ## Clean build artifacts and test cache
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@$(GO_CMD) clean
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean complete!$(NC)"

fmt: ## Format Go code
	@echo "$(YELLOW)Formatting Go code...$(NC)"
	@$(GO_CMD) fmt ./...
	@echo "$(GREEN)Code formatted!$(NC)"

lint: ## Run linter
	@echo "$(YELLOW)Running linter...$(NC)"
	@golangci-lint run
	@echo "$(GREEN)Linting complete!$(NC)"

deps: ## Download and tidy dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	@$(GO_CMD) mod download
	@$(GO_CMD) mod tidy
	@echo "$(GREEN)Dependencies updated!$(NC)"

# Install migrate tool if not present
install-migrate: ## Install golang-migrate tool
	@echo "$(YELLOW)Installing golang-migrate...$(NC)"
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(GREEN)golang-migrate installed!$(NC)"

# Check if required tools are installed
check-tools: ## Check if required tools are installed
	@echo "$(YELLOW)Checking required tools...$(NC)"
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)Docker is required but not installed.$(NC)"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "$(RED)Docker Compose is required but not installed.$(NC)"; exit 1; }
	@command -v migrate >/dev/null 2>&1 || { echo "$(YELLOW)golang-migrate not found. Run 'make install-migrate' to install it.$(NC)"; }
	@echo "$(GREEN)Tool check complete!$(NC)"

# Full setup for new developers
init: check-tools setup db-up migrate-up ## Initialize project for new developers
	@echo "$(GREEN)Project initialization complete!$(NC)"
	@echo "$(BLUE)You can now run 'make run' to start the application$(NC)"
