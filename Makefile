.PHONY: help migrate-up migrate-down migrate-status migrate-create

help:
	@echo "Available commands:"
	@echo "  make migrate-up        - Apply all pending migrations"
	@echo "  make migrate-down      - Rollback last migration"
	@echo "  make migrate-status    - Show migration status"
	@echo "  make migrate-create    - Create new migration (use NAME=migration_name)"

# Database connection from config
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASS ?= root
DB_NAME ?= crm
DB_SSLMODE ?= disable

DB_URL := "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

migrate-up:
	@echo "Applying migrations..."
	goose -dir migrations postgres $(DB_URL) up

migrate-down:
	@echo "Rolling back last migration..."
	goose -dir migrations postgres $(DB_URL) down

migrate-status:
	@echo "Migration status:"
	goose -dir migrations postgres $(DB_URL) status

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=your_migration_name"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(NAME)"
	goose -dir migrations create $(NAME) sql