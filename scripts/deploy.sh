#!/bin/bash
set -euo pipefail

ENV=${1:-production}
ROOT_DIR="$(cd "$(dirname "$0")/.. && pwd)"
echo "Deploy mode: $ENV"
echo "Project dir: $ROOT_DIR"

cd "$ROOT_DIR"

if [ ! -f backend/.env ]; then
  if [ -f backend/.env.example ]; then
    cp backend/.env.example backend/.env
  else
    echo "Error: backend/.env not found"
    exit 1
  fi
fi

echo "Running database migrations..."
cd backend && go run cmd/migrate/main.go
cd "$ROOT_DIR"

echo "Building backend binary..."
cd backend && CGO_ENABLED=0 go build -o server ./cmd/server/
cd "$ROOT_DIR"

if [ -d pc ]; then
  echo "Building PC frontend..."
  cd pc && npm ci --silent && npx vite build
  cd "$ROOT_DIR"
fi

if [ -d admin ]; then
  echo "Building admin..."
  cd admin && npm ci --silent && npx vite build
  cd "$ROOT_DIR"
fi

if [ "$ENV" = "production" ]; then
  echo "Starting Docker Compose..."
  docker compose up -d --build
  echo "Services:"
  docker compose ps
  echo "Deploy complete!"
else
  echo "Starting backend..."
  cd backend && nohup ./server > /tmp/server.log 2>&1 &
  echo "Backend PID: $!"
  cd "$ROOT_DIR"
  echo "Local deploy ready."
  echo "  Backend: http://localhost:8080"
fi
