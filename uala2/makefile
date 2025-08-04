SHELL := /bin/bash

# Directorios y binarios
ROOT_DIR := $(shell pwd)
BUILD_DIR := $(ROOT_DIR)/bin
BINARY := $(BUILD_DIR)/bootstrap

# Docker Compose y SAM CLI
DC := docker compose -f docker-compose.yml
SAM := sam

.PHONY: build compose-up sam-up start restart down clean

# 1. Compilar la Lambda Go
build:
	@echo "🔨 Compiling Lambda (Go bootstrap)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -gcflags "all=-N -l" \
		-o $(BINARY) cmd/main.go
	@chmod +x $(BINARY)
	@echo "✅ Build complete"

# 2. Levantar servicios Docker (DynamoDB, Lambda Go Dev y LocalStack)
compose-up:
	@echo "⬆️  Starting Docker services..."
	$(DC) up -d
	@echo "✅ Docker services are up"

# 3. Iniciar SAM CLI (API Gateway → Lambda)
sam-up:
	@echo "🚀 Starting SAM local API..."
	$(SAM) local start-api \
		--host 0.0.0.0 \
		--port 3000 \
		--docker-network app-network
	@echo "✅ SAM local API running on http://localhost:3000"

# 4. Shortcut: build + compose-up + sam-up
start: build compose-up sam-up
	@echo "🏁 All set! Development environment ready"

# 5. Tear down servicios Docker y binarios
stop:
	@echo "⬇️  Stopping Docker services..."
	$(DC) down --remove-orphans
	@echo "✅ Docker services stopped"

# 6. Reiniciar todo desde cero (down + start)
restart: stop clean start
	@echo "🔄 Restart complete"

# 7. Borrar binarios compilados
clean:
	@echo "🧹 Performing full clean..."
	$(DC) down -v --remove-orphans
	@echo "🗑️  Removed Docker volumes and networks"
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	@echo "✅ Full clean complete"


