.PHONY: build install clean test docker-build docker-push help

# Binary name
BINARY_NAME=turso-migrate

# Docker image name
DOCKER_IMAGE=turso-migrate
DOCKER_REGISTRY=ghcr.io/rubenmeza

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build targets
build: ## Build the CLI binary
	@echo "Building $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BINARY_NAME) ./cmd/turso-migrate

install: ## Install the CLI binary using go install
	@echo "Installing $(BINARY_NAME)..."
	@$(GOCMD) install ./cmd/turso-migrate

clean: ## Remove binary and clean up
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)

test: ## Run tests
	@echo "Running tests..."
	@$(GOTEST) -v ./...

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@$(GOMOD) tidy

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest..."
	@docker build -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest .
	@docker tag $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest $(DOCKER_IMAGE):latest

docker-push: docker-build ## Push Docker image to registry
	@echo "Pushing Docker image $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest..."
	@docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest

docker-run: ## Run Docker image locally (requires TURSO_DATABASE_URL and TURSO_AUTH_TOKEN env vars)
	@echo "Running Docker image locally..."
	@docker run --rm -e TURSO_DATABASE_URL -e TURSO_AUTH_TOKEN $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest

# Development targets
dev-deps: ## Install development dependencies
	@echo "Installing development dependencies..."
	@$(GOGET) -t ./...

format: ## Format code using gofmt
	@echo "Formatting code..."
	@gofmt -s -w .

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run

# Release targets
tag: ## Create a git tag (usage: make tag VERSION=v1.0.0)
	@echo "Creating tag $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)

release: clean test build ## Prepare for release
	@echo "Release build completed!"

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Default target
.DEFAULT_GOAL := help