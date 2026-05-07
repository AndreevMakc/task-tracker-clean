.PHONY: help migrate run check-migrations

# Load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Default values (can be overridden by .env or environment)
PG_URL ?= postgres://postgres:postgres@localhost:5432/tasktracker?sslmode=disable
MIGRATIONS_PATH ?= migrations

help:
	@echo "Available commands:"
	@echo "  make migrate          - Apply all pending database migrations"
	@echo "  make run              - Apply migrations and start the application"
	@echo "  make check-migrations - Show current migration version in database"
	@echo "  make help             - Show this help message"

migrate:
	@echo "Applying migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(PG_URL)" up

check-migrations:
	@echo "Current migration version:"
	migrate -path $(MIGRATIONS_PATH) -database "$(PG_URL)" version

run: migrate
	@echo "Starting application..."
	go run ./cmd/app
