.PHONY: build run test clean fmt vet docker-build docker-run help

# Build the application
build:
	@echo "Building telegram-sync..."
	@go build -o telegram-sync .

# Run the application
run:
	@echo "Running telegram-sync..."
	@go run main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f telegram-sync
	@rm -rf downloads/

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run linters
lint: fmt vet
	@echo "Linting complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t telegram-sync:latest .

# Run with Docker Compose
docker-run:
	@echo "Starting with Docker Compose..."
	@docker-compose up -d

# Stop Docker Compose
docker-stop:
	@echo "Stopping Docker Compose..."
	@docker-compose down

# View Docker logs
docker-logs:
	@docker-compose logs -f

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# Help command
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make fmt           - Format code"
	@echo "  make vet           - Run go vet"
	@echo "  make lint          - Run all linters"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run with Docker Compose"
	@echo "  make docker-stop   - Stop Docker Compose"
	@echo "  make docker-logs   - View Docker logs"
	@echo "  make deps          - Install dependencies"
	@echo "  make deps-update   - Update dependencies"
	@echo "  make help          - Show this help message"
