#!/bin/bash

# Migration script for HeiCRM
# Usage: ./scripts/migrate.sh [up|down|status|create]

set -e

# Default database configuration (can be overridden by environment variables)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASS="${DB_PASS:-root}"
DB_NAME="${DB_NAME:-crm}"
DB_SSLMODE="${DB_SSLMODE:-disable}"

DB_URL="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

MIGRATIONS_DIR="migrations"

# Check if goose is installed
if ! command -v goose &> /dev/null; then
    echo "Error: goose is not installed"
    echo "Install it with: go install github.com/pressly/goose/v3/cmd/goose@latest"
    exit 1
fi

case "$1" in
    up)
        echo "Applying all pending migrations..."
        goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" up
        ;;
    down)
        echo "Rolling back last migration..."
        goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" down
        ;;
    status)
        echo "Migration status:"
        goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" status
        ;;
    create)
        if [ -z "$2" ]; then
            echo "Error: Migration name is required"
            echo "Usage: $0 create <migration_name>"
            exit 1
        fi
        echo "Creating new migration: $2"
        goose -dir "$MIGRATIONS_DIR" create "$2" sql
        ;;
    *)
        echo "Usage: $0 {up|down|status|create <name>}"
        echo ""
        echo "Commands:"
        echo "  up              - Apply all pending migrations"
        echo "  down            - Rollback last migration"
        echo "  status          - Show migration status"
        echo "  create <name>   - Create new migration file"
        echo ""
        echo "Environment variables:"
        echo "  DB_HOST     (default: localhost)"
        echo "  DB_PORT     (default: 5432)"
        echo "  DB_USER     (default: postgres)"
        echo "  DB_PASS     (default: root)"
        echo "  DB_NAME     (default: crm)"
        echo "  DB_SSLMODE  (default: disable)"
        exit 1
        ;;
esac
