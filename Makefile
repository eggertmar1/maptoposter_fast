.PHONY: all build clean run test docker-build docker-run help

# Variables
BINARY_SERVER=server
BINARY_CLI=maptoposter-cli
DOCKER_IMAGE=maptoposter-fast

all: build

# Build both server and CLI
build:
	@echo "Building server..."
	@go build -o $(BINARY_SERVER) ./cmd/server
	@echo "Building CLI..."
	@go build -o $(BINARY_CLI) ./cmd/cli
	@echo "✓ Build complete!"

# Build only the server
server:
	@echo "Building server..."
	@go build -o $(BINARY_SERVER) ./cmd/server
	@echo "✓ Server built!"

# Build only the CLI
cli:
	@echo "Building CLI..."
	@go build -o $(BINARY_CLI) ./cmd/cli
	@echo "✓ CLI built!"

# Run the server
run: server
	@echo "Starting server..."
	@./$(BINARY_SERVER)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_SERVER) $(BINARY_CLI)
	@rm -rf posters/*.png
	@echo "✓ Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies installed!"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "✓ Docker image built!"

# Run with Docker Compose
docker-run:
	@echo "Starting service with Docker Compose..."
	@docker-compose up -d
	@echo "✓ Service started at http://localhost:8080"

# Stop Docker Compose
docker-stop:
	@echo "Stopping service..."
	@docker-compose down
	@echo "✓ Service stopped!"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted!"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Show help
help:
	@echo "Map Poster Fast - Makefile commands:"
	@echo ""
	@echo "  make build        - Build both server and CLI"
	@echo "  make server       - Build only the server"
	@echo "  make cli          - Build only the CLI"
	@echo "  make run          - Build and run the server"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make deps         - Install dependencies"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run with Docker Compose"
	@echo "  make docker-stop  - Stop Docker Compose"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"
	@echo "  make help         - Show this help"
