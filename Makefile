.PHONY: help setup build up down logs migrate seed test test-backend test-frontend lint clean reset

help:
	@echo "Task Management App - Available Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make setup          - Initial setup (copy .env files)"
	@echo ""
	@echo "Development:"
	@echo "  make up             - Start all services (Docker Compose)"
	@echo "  make down           - Stop all services"
	@echo "  make logs           - View logs from all services"
	@echo ""
	@echo "Database:"
	@echo "  make migrate        - Run database migrations"
	@echo "  make seed           - Seed database with sample data"
	@echo ""
	@echo "Build:"
	@echo "  make build          - Build backend and frontend"
	@echo "  make build-backend  - Build backend only"
	@echo "  make build-frontend - Build frontend only"
	@echo ""
	@echo "Testing:"
	@echo "  make test           - Run all tests"
	@echo "  make test-backend   - Run backend tests"
	@echo "  make test-frontend  - Run frontend tests"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint           - Lint all code"
	@echo "  make lint-backend   - Lint backend only"
	@echo "  make lint-frontend  - Lint frontend only"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make reset          - Reset everything (Docker, volumes)"

setup:
	@if [ ! -f .env ]; then cp .env.example .env; echo "✓ .env created"; else echo ".env already exists"; fi
	@if [ ! -f backend/.env ]; then cp .env.example backend/.env; echo "✓ backend/.env created"; fi
	@if [ ! -f frontend/.env.local ]; then echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > frontend/.env.local; echo "✓ frontend/.env.local created"; fi

up: setup
	docker-compose up -d
	@echo "✓ Services started. Backend: http://localhost:8080, Frontend: http://localhost:3000"

down:
	docker-compose down
	@echo "✓ Services stopped"

logs:
	docker-compose logs -f

build: build-backend build-frontend

build-backend:
	docker-compose build backend

build-frontend:
	docker-compose build frontend

migrate:
	docker-compose exec -T backend /app/taskapp migrate

seed:
	docker-compose exec -T backend /app/taskapp seed

test: test-backend test-frontend

test-backend:
	docker run --rm -v "$(PWD)/backend":/app -w /app golang:1.21-alpine go test -v ./...

test-frontend:
	docker run --rm -v "$(PWD)/frontend":/app -w /app node:20-alpine sh -c "npm install --silent && npm test -- --watchAll=false --forceExit"

lint: lint-backend lint-frontend

lint-backend:
	docker run --rm -v "$(PWD)/backend":/app -w /app golang:1.21-alpine sh -c "gofmt -w . && go vet ./..."

lint-frontend:
	cd frontend && npm run lint

clean:
	rm -rf backend/bin
	rm -rf frontend/.next
	rm -rf frontend/node_modules
	docker-compose down

reset: clean
	docker-compose down -v
	@echo "✓ Everything cleaned and reset"

.DEFAULT_GOAL := help
