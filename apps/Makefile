.PHONY: help db-up db-down db-logs backend-up backend-down backend-build frontend-up frontend-down frontend-build test-e2e test-e2e-ui all-up all-down

BACKEND_DIR := backend
FRONTEND_DIR := frontend

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Database:"
	@echo "  db-up       Start PostgreSQL container"
	@echo "  db-down     Stop PostgreSQL container"
	@echo "  db-logs     Show PostgreSQL logs"
	@echo ""
	@echo "Backend (Go):"
	@echo "  backend-up      Start backend with hot reload (air)"
	@echo "  backend-down    Stop backend process"
	@echo "  backend-build   Build backend binary"
	@echo ""
	@echo "Frontend (Next.js):"
	@echo "  frontend-up      Start frontend dev server"
	@echo "  frontend-down    Stop frontend process"
	@echo "  frontend-build   Build frontend for production"
	@echo ""
	@echo "Tests (Playwright):"
	@echo "  test-e2e       Run E2E tests in headless mode"
	@echo "  test-e2e-ui    Run E2E tests with browser UI (headed)"
	@echo ""
	@echo "All:"
	@echo "  all-up      Start database, backend, and frontend"
	@echo "  all-down    Stop everything"

# Database
db-up:
	docker compose up -d

db-down:
	docker compose down

db-logs:
	docker compose logs -f postgres

# Backend
backend-up:
	cd $(BACKEND_DIR) && air

backend-down:
	@pkill -f "air" 2>/dev/null || true

backend-build:
	cd $(BACKEND_DIR) && go build ./cmd/server

# Frontend
frontend-up:
	cd $(FRONTEND_DIR) && npm run dev

frontend-down:
	@pkill -f "next dev" 2>/dev/null || true

frontend-build:
	cd $(FRONTEND_DIR) && npm run build

# Tests (Playwright)
test-e2e:
	cd $(FRONTEND_DIR) && npx playwright test

test-e2e-ui:
	cd $(FRONTEND_DIR) && npx playwright test --headed

# All
all-up: db-up
	@echo "Starting backend in background..."
	@cd $(BACKEND_DIR) && air &
	@sleep 2
	@echo "Starting frontend in background..."
	@cd $(FRONTEND_DIR) && npm run dev &
	@echo "All services started."

all-down: backend-down frontend-down db-down
	@echo "All services stopped."
