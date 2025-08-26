.PHONY: help build run clean docker-up docker-down docker-logs deps migrate migrate-reset migrate-status

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
	@echo "  migrate       - Run database migrations"
	@echo "  migrate-reset - Reset database and run all migrations"
	@echo "  migrate-status- Check migration status"

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


# Database commands
migrate: ## Run database migrations
	@echo "🗄️  Running database migrations..."
	@for file in migrations/*.sql; do \
		echo "Running $$(basename $$file)..."; \
		docker exec huddle_postgres psql -U huddle_user -d huddle -f /docker-entrypoint-initdb.d/$$(basename $$file); \
	done
	@echo "✅ All migrations completed successfully"

migrate-reset: ## Reset database and run all migrations
	@echo "🔄 Resetting database and running migrations..."
	@make docker-down
	@make docker-up
	@echo "✅ Database reset and migrations completed"

migrate-status: ## Check migration status
	@echo "📊 Checking migration status..."
	@docker exec huddle_postgres psql -U huddle_user -d huddle -c "\dt"
	@echo "✅ Migration status checked"

# Development helpers
dev: docker-up deps run

restart: docker-down docker-up
	@echo "🔄 Services restarted"
