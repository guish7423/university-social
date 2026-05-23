#!/bin/bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.. && pwd)"
cd "$ROOT_DIR/backend"

echo "Running migrations..."
go run cmd/migrate/main.go

echo "Seeding data..."
go run cmd/seed/main.go --count 30

echo "Database ready!"
