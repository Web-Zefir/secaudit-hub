.PHONY: up down dev dev-backend dev-frontend build test

up:
	docker compose -f docker-compose.dev.yml up -d
	docker compose -f docker-compose.dev.yml ps

down:
	docker compose -f docker-compose.dev.yml down

dev:
	$(MAKE) -j2 dev-backend dev-frontend

dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && npm run dev

build:
	cd backend && go build -o bin/server ./cmd/server
	cd frontend && npm run build

test:
	cd backend && go test ./...
	cd frontend && npm run build
