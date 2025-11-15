.PHONY: build up down logs restart clean db-logs migrate

# Build and start all services
build:
	docker-compose build --no-cache

# Start all services
up:
	docker-compose up -d

# Start with logs
up-logs:
	docker-compose up

# Stop all services
down:
	docker-compose down

# Show logs
logs:
	docker-compose logs -f

# Show app logs only
app-logs:
	docker-compose logs -f app

# Show database logs only
db-logs:
	docker-compose logs -f postgres

# Restart all services
restart:
	docker-compose restart

# Clean up everything
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f

# Run database migrations (if using migrate tool)
migrate-up:
	docker-compose exec app migrate -path /app/migrations -database "postgres://goevent_user:goevent_password@postgres:5432/goevent?sslmode=disable" up

# Create new migration
migrate-create:
	@echo "Usage: make migrate-create name=migration_name"
	docker-compose exec app migrate create -ext sql -dir /app/migrations -seq $(name)

# Access database shell
db-shell:
	docker-compose exec postgres psql -U goevent_user -d goevent

# Run tests (if any)
test:
	docker-compose exec app go test ./...

# Build and run locally (without Docker)
run-local:
	go run cmd/api/main.go

# Install dependencies locally
deps:
	go mod tidy
	go mod download