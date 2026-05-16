.PHONY: up down logs backend frontend migrate-up migrate-down test fmt help

# Load .env if present
ifneq (,$(wildcard .env))
  include .env
  export
endif

DB_URL ?= postgres://treinozap:treinozap@localhost:5432/treinozap?sslmode=disable

help:
	@echo "TreinoZap — comandos disponíveis:"
	@echo "  make up            Sobe todos os serviços com Docker Compose"
	@echo "  make down          Para e remove os containers"
	@echo "  make logs          Exibe logs de todos os serviços"
	@echo "  make backend       Sobe apenas o backend"
	@echo "  make frontend      Sobe apenas o frontend"
	@echo "  make migrate-up    Aplica todas as migrations"
	@echo "  make migrate-down  Reverte a última migration"
	@echo "  make test          Roda os testes do backend"
	@echo "  make fmt           Formata o código Go"

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

backend:
	docker compose up -d backend

frontend:
	docker compose up -d frontend

migrate-up:
	@echo "Aplicando migrations..."
	cd backend && go run ./cmd/migrate/main.go up

migrate-down:
	@echo "Revertendo última migration..."
	cd backend && go run ./cmd/migrate/main.go down

test:
	cd backend && go test ./...

fmt:
	cd backend && gofmt -w .
	cd frontend && npm run lint --fix 2>/dev/null || true

dev-backend:
	cd backend && go run ./cmd/api/main.go

dev-frontend:
	cd frontend && npm run dev
