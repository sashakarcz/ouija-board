.PHONY: build run clean test docker-build docker-run docker-stop lint fmt vet help

# Variables
BINARY_NAME=ouija-board
DOCKER_IMAGE=ouija-board:go
DOCKER_CONTAINER=ouija-board-container

# Default target
.DEFAULT_GOAL := help

## build: Build the Go binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

## run: Build and run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f answers.json
	@rm -f answers.json.tmp
	@go clean
	@echo "Clean complete"

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "Tests complete"

## test-coverage: Run tests with coverage report
test-coverage: test
	@go tool cover -html=coverage.out

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	@docker build -f Dockerfile-go -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -d --name $(DOCKER_CONTAINER) -p 8080:8080 \
		-e OLLAMA_URL="http://host.docker.internal:11434/api/generate" \
		$(DOCKER_IMAGE)
	@echo "Container started: $(DOCKER_CONTAINER)"

## docker-stop: Stop and remove Docker container
docker-stop:
	@echo "Stopping Docker container..."
	@docker stop $(DOCKER_CONTAINER) || true
	@docker rm $(DOCKER_CONTAINER) || true
	@echo "Container stopped and removed"

## docker-logs: Show Docker container logs
docker-logs:
	@docker logs -f $(DOCKER_CONTAINER)

## compose-up: Start with docker-compose
compose-up:
	@echo "Starting with docker-compose..."
	@docker-compose -f docker-compose-go.yaml up -d
	@echo "Services started"

## compose-down: Stop docker-compose services
compose-down:
	@echo "Stopping docker-compose services..."
	@docker-compose -f docker-compose-go.yaml down
	@echo "Services stopped"

## compose-logs: Show docker-compose logs
compose-logs:
	@docker-compose -f docker-compose-go.yaml logs -f

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install from https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run
	@echo "Linting complete"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Formatting complete"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies downloaded"

## security-check: Check for security vulnerabilities
security-check:
	@echo "Checking for security vulnerabilities..."
	@which gosec > /dev/null || (echo "gosec not found. Install with: go install github.com/securego/gosec/v2/cmd/gosec@latest" && exit 1)
	@gosec ./...
	@echo "Security check complete"

## all: Clean, format, vet, test, and build
all: clean fmt vet test build
	@echo "All tasks complete"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
