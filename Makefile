.PHONY: help up down start stop restart logs build clean test install

help: ## Show this help message
	@echo "ChatMeNow - Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

up: ## Start all containers (alias for start)
	@echo "Starting ChatMeNow containers..."
	@docker compose up -d
	@echo "Containers started!"

down: ## Stop and remove all containers (alias for stop)
	@echo "Stopping and removing ChatMeNow containers..."
	@docker compose down
	@echo "Containers stopped and removed!"

start: ## Start all services
	@echo "Starting ChatMeNow..."
	@docker compose up -d
	@echo "Services started!"
	
stop: ## Stop all services
	@echo "Stopping ChatMeNow..."
	@docker compose down
	@echo "Services stopped!"

restart: ## Restart all services
	@echo "Restarting ChatMeNow..."
	@docker compose restart
	@echo "Services restarted!"

logs: ## Show logs (use SERVICE=gateway to see specific service)
ifdef SERVICE
	@docker compose logs -f $(SERVICE)
else
	@docker compose logs -f
endif

build: ## Rebuild all services
	@echo "Building services..."
	@docker compose build
	@echo "Build complete!"

clean: ## Remove all containers, volumes, and networks
	@echo "Cleaning up..."
	@docker compose down -v
	@echo "Cleanup complete!"

test: ## Run API tests
	@echo "Running tests..."
	@bash test-api.sh

install-gateway: ## Install gateway dependencies
	@cd gateway && npm install

install-auth: ## Install auth-service dependencies
	@cd auth-service && npm install

install-blog: ## Install blog-service dependencies
	@cd blog-service && npm install

install-chat: ## Install chat-service dependencies
	@cd chat-service && go mod download

install: install-gateway install-auth install-blog install-chat ## Install all dependencies

dev-gateway: ## Run gateway in development mode
	@cd gateway && npm run start:dev

dev-auth: ## Run auth-service in development mode
	@cd auth-service && npm run start:dev

dev-blog: ## Run blog-service in development mode
	@cd blog-service && npm run start:dev

dev-chat: ## Run chat-service in development mode
	@cd chat-service && go run cmd/server/main.go

db-postgres: ## Connect to PostgreSQL
	@docker exec -it chatmenow-postgres psql -U chatmenow -d chatmenow

db-mongo: ## Connect to MongoDB
	@docker exec -it chatmenow-mongodb mongosh -u chatmenow -p chatmenow123

db-redis: ## Connect to Redis
	@docker exec -it chatmenow-redis redis-cli

backup-db: ## Backup databases
	@echo "Backing up databases..."
	@mkdir -p backups
	@docker exec chatmenow-postgres pg_dump -U chatmenow chatmenow > backups/postgres_$(shell date +%Y%m%d_%H%M%S).sql
	@docker exec chatmenow-mongodb mongodump --username=chatmenow --password=chatmenow123 --out=/tmp/backup
	@docker cp chatmenow-mongodb:/tmp/backup backups/mongodb_$(shell date +%Y%m%d_%H%M%S)
	@echo "Backup complete! Check ./backups/"

status: ## Show service status
	@docker compose ps
