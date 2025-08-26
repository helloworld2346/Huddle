.PHONY: help build run clean docker-up docker-down docker-logs deps

# Default target
help:
	@echo "Available commands:"
	@echo "  docker-up     - Start PostgreSQL and Redis containers"
	@echo "  docker-down   - Stop all containers"
	@echo "  docker-logs   - Show container logs"
	@echo "  deps          - Download Go dependencies"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  clean         - Clean build artifacts"

# Docker commands
docker-up:
	@echo "🚀 Starting PostgreSQL and Redis containers..."
	docker-compose up -d
	@echo "✅ Containers started successfully"
	@echo "📊 PostgreSQL: localhost:5432"
	@echo "🔴 Redis: localhost:6379"

docker-down:
	@echo "🛑 Stopping containers..."
	docker-compose down
	@echo "✅ Containers stopped"

docker-logs:
	docker-compose logs -f

# Go commands
deps:
	@echo "📦 Downloading dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies downloaded"

build:
	@echo "🔨 Building application..."
	go build -o bin/huddle ./cmd/server
	@echo "✅ Application built successfully"

run:
	@echo "🚀 Running application..."
	go run ./cmd/server

clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	go clean
	@echo "✅ Cleaned successfully"

# Development helpers
dev: docker-up deps run

restart: docker-down docker-up
	@echo "🔄 Services restarted"
