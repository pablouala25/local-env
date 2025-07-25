SHELL := /bin/bash

ROOT_DIR := $(shell pwd)
BUILD_DIR := $(ROOT_DIR)/bin
VERSION := 1.0
DOCKER_COMPOSE_FILE := $(ROOT_DIR)/docker-compose.yml

# >>> Cambiado: el binario Lambda debe llamarse *bootstrap*
BINARY_NAME := bootstrap
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)

CMD_PATH := $(ROOT_DIR)/cmd/main.go

# Detectar arch opcional (amd64 default)
GOARCH ?= amd64

.PHONY: build run test clean lint up down logs compose-build dev restart reload invoke

build:
	@echo "🔨 Building Go Lambda bootstrap ($(GOARCH))..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) \
		go build -gcflags "all=-N -l" \
		-o $(BINARY_PATH) \
		-ldflags "-X main.Version=$(VERSION)" \
		$(CMD_PATH)
	@chmod +x $(BINARY_PATH)

run:
	@echo "🚀 Running app (local direct exec)..."
	@go run $(CMD_PATH)

test:
	@echo "🧪 Running tests..."
	@go test ./...

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

lint:
	@echo "🔍 Running linter..."
	@golangci-lint run --config .golangci.yml --verbose

up:
	@echo "⬆️  Starting Docker Compose services..."
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d
	@$(MAKE) logs

down:
	@echo "⬇️  Stopping Docker Compose services..."
	docker compose -f $(DOCKER_COMPOSE_FILE) down --remove-orphans

logs:
	@echo "📋 Tailing Docker Compose logs..."
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

compose-build:
	@echo "🐳 Building Docker Compose images..."
	docker compose -f $(DOCKER_COMPOSE_FILE) build

dev: build up
	@echo "✅ Dev environment is up and ready!"

restart:
	@echo "🔄 Restarting lambda-go-dev (to reload updated binary)..."
	docker compose -f $(DOCKER_COMPOSE_FILE) restart lambda-go-dev

reload: build restart
	@echo "♻️  Binario reconstruido y lambda-go-dev reiniciado."
