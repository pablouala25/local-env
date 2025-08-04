SHELL := /bin/bash

# Directories and binaries
ROOT_DIR := $(shell pwd)
BUILD_DIR := $(ROOT_DIR)/bin
BINARY := $(BUILD_DIR)/bootstrap

# Docker Compose and SAM CLI
DC := docker compose -f docker-compose.yml
SAM := sam

.PHONY: build compose-up sam-up restart-all down clean all

# 1. Compile the Go Lambda
build:
	@echo "🔨 Compiling Lambda (Go bootstrap)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -gcflags "all=-N -l" \
		-o $(BINARY) cmd/main.go
	@chmod +x $(BINARY)

# 2. Bring up external services and the Go container
compose-up:
	@echo "⬆️  Starting external services..."
	$(DC) up -d

# 3. Start SAM local for API Gateway / Lambda emulation
sam-up:
	@echo "🚀 Starting SAM local API..."
	$(SAM) local start-api \
		--host 0.0.0.0 \
		--port 3000 \
		--docker-network app-network

# 4. Rebuild & restart the full dev stack (Go code + containers)
restart-all:
	@echo "🔄 Restarting full dev stack..."
	$(MAKE) build
	$(DC) restart
	@# SAM local hot-reloads function code automatically; no restart needed

# 5. Tear down all services and clean up binaries
down:
	@echo "⬇️  Stopping and removing containers, networks, volumes..."
	$(DC) down --remove-orphans
	rm -rf $(BUILD_DIR)

# 6. One-step shortcut: compile, bring up services, and start SAM
all: build compose-up sam-up
	@echo "✅ Development environment ready! Visit http://localhost:3000"
