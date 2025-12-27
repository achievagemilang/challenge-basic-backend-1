.PHONY: run run-worker test dev docker-up docker-down docker-down-v swag migrate-up migrate-down migrate-reset migrate-create test-coverage

# Run the web application
run:
	go run cmd/web/main.go

# Run the worker application
run-worker:
	go run cmd/worker/main.go

# Run unit tests
test:
	go test -v ./test/...

test-coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

# Run with Hot Reload (Air)
dev:
	air

# Docker commands
docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-down-v:
	docker compose down -v

# Swagger
swag:
	swag init -g cmd/web/main.go --parseDependency --parseInternal

# Database Migrations
DB_URL := "postgres://postgres:postgres@localhost:54320/challenge_backend_db?sslmode=disable"
MIGRATE_CMD := migrate -database $(DB_URL) -path db/migrations

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

migrate-reset:
	$(MIGRATE_CMD) reset

# Usage: make migrate-create name=create_table_users
migrate-create:
	migrate create -ext sql -dir db/migrations $(name)
