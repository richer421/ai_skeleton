.PHONY: help frontend-dev backend-dev gen-swagger gen-sql

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

frontend-dev: ## Run frontend in development mode
	cd frontend && npm run dev

backend-dev: ## Run backend in development mode with hot reload
	cd backend && air

gen-swagger: ## Generate Swagger documentation
	cd backend && swag init -g cmd/server/main.go

gen-sql: ## Generate SQL code with gen-gorm
	cd backend && go run cmd/gen/main.go
