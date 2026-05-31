.PHONY: fmt vet lint test coverage build migrate migrate-down seed security secrets check docker-up docker-down dev-gateway dev-api dev-worker

GO ?= go
GOLANGCI_LINT ?= golangci-lint
GOFUMPT ?= gofumpt
MIGRATE ?= migrate

BIN_DIR := bin
GATEWAY_BIN := $(BIN_DIR)/gateway
API_BIN := $(BIN_DIR)/api
WORKER_BIN := $(BIN_DIR)/worker

fmt:
	$(GO) fmt ./...
	$(GOFUMPT) -w .

vet:
	$(GO) vet ./...

lint:
	$(GOLANGCI_LINT) run ./...

test:
	$(GO) test ./... -race -count=1

coverage:
	$(GO) test ./... -race -coverprofile=coverage.out $(shell go list ./... | grep -v /cmd/)
	@$(GO) tool cover -func=coverage.out | tail -1

build: $(GATEWAY_BIN) $(API_BIN) $(WORKER_BIN)

$(GATEWAY_BIN):
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(GATEWAY_BIN) ./cmd/gateway

$(API_BIN):
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(API_BIN) ./cmd/api

$(WORKER_BIN):
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(WORKER_BIN) ./cmd/worker

migrate:
	$(MIGRATE) -path migrations -database "$(BASALT_DATABASE_DSN)" up

migrate-down:
	$(MIGRATE) -path migrations -database "$(BASALT_DATABASE_DSN)" down 1

seed:
	$(GO) run ./cmd/seed

security:
	govulncheck ./...

secrets:
	gitleaks detect --no-git --source .

check: fmt vet lint test build

docker-up:
	docker compose -f deploy/docker-compose.yml up -d --build

docker-down:
	docker compose -f deploy/docker-compose.yml down

dev-gateway:
	air -c .air.gateway.toml

dev-api:
	air -c .air.api.toml

dev-worker:
	air -c .air.worker.toml
