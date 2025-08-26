# Makefile for knm (KryptoNim) project

# Variables
BINARY_NAME=kryptonim-app
BUILD_DIR=build
DOCKER_IMAGE=knm
DOCKER_TAG=latest

# Default target
.DEFAULT_GOAL := help

.PHONY: build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/kryptonim

.PHONY: run
run: ## Run the application (requires OPENEXCHANGERATES_APP_ID env var)
	@echo "Running $(BINARY_NAME)..."
	go run ./cmd/kryptonim

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	go clean
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 -e OPENEXCHANGERATES_APP_ID=your_app_id_here $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-clean
docker-clean: ## Clean Docker images
	@echo "Cleaning Docker images..."
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || true
