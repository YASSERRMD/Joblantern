# Joblantern Makefile
# Run `make help` for the list of targets.

SHELL := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

GO              ?= go
GOFLAGS         ?=
PKG             := ./...
BIN_DIR         := bin
BIN             := $(BIN_DIR)/joblantern
COMPOSE         ?= docker compose
COMPOSE_FILE    := deploy/docker-compose.yml
GOLANGCI_LINT   ?= golangci-lint
GOOSE           ?= goose
SQLC            ?= sqlc
TEMPL           ?= templ

DATABASE_URL    ?= postgres://joblantern:joblantern@localhost:5432/joblantern?sslmode=disable

# ---------------------------------------------------------------------------
# Help
# ---------------------------------------------------------------------------
.PHONY: help
help: ## Print this help.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make <target>\n\nTargets:\n"} \
	      /^[a-zA-Z_-]+:.*?##/ {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' \
	      $(MAKEFILE_LIST)

# ---------------------------------------------------------------------------
# Build / run
# ---------------------------------------------------------------------------
.PHONY: build
build: ## Build the joblantern binary into bin/.
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN) ./cmd/joblantern

.PHONY: run
run: build ## Build and run the joblantern server.
	$(BIN)

# ---------------------------------------------------------------------------
# Quality gates
# ---------------------------------------------------------------------------
.PHONY: fmt
fmt: ## Format all Go code.
	$(GO) fmt $(PKG)

.PHONY: lint
lint: ## Run golangci-lint (skipped silently if not installed).
	@if command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
	    $(GOLANGCI_LINT) run $(PKG); \
	else \
	    echo "golangci-lint not installed; skipping (install: https://golangci-lint.run/usage/install/)"; \
	fi

.PHONY: test
test: ## Run unit tests with race detector.
	$(GO) test $(GOFLAGS) -race -count=1 $(PKG)

.PHONY: test-integration
test-integration: ## Run integration tests (build tag: integration).
	$(GO) test $(GOFLAGS) -race -count=1 -tags=integration $(PKG)

.PHONY: vet
vet: ## Run go vet.
	$(GO) vet $(PKG)

.PHONY: tidy
tidy: ## Run go mod tidy.
	$(GO) mod tidy

# ---------------------------------------------------------------------------
# Database / codegen (placeholders, wired in Phase 02+)
# ---------------------------------------------------------------------------
.PHONY: migrate-up
migrate-up: ## Apply all pending migrations.
	$(GO) run ./cmd/goose-migrate -dsn "$(DATABASE_URL)" -dir migrations up

.PHONY: migrate-down
migrate-down: ## Roll back the most recent migration.
	$(GO) run ./cmd/goose-migrate -dsn "$(DATABASE_URL)" -dir migrations down

.PHONY: migrate-status
migrate-status: ## Show migration status.
	$(GO) run ./cmd/goose-migrate -dsn "$(DATABASE_URL)" -dir migrations status

.PHONY: migrate-reset
migrate-reset: ## Roll all migrations back to zero (development only).
	$(GO) run ./cmd/goose-migrate -dsn "$(DATABASE_URL)" -dir migrations reset

.PHONY: migrate-create
migrate-create: ## Create a new SQL migration: NAME=add_thing make migrate-create
	@test -n "$(NAME)" || (echo "NAME is required (e.g. NAME=add_thing make migrate-create)"; exit 1)
	$(GO) run ./cmd/goose-migrate -dir migrations create $(NAME) sql

.PHONY: sqlc
sqlc: ## Generate Go from SQL queries.
	$(SQLC) generate

.PHONY: templ
templ: ## Generate Go from templ templates.
	$(TEMPL) generate

.PHONY: tailwind
tailwind: ## Build Tailwind CSS bundle.
	@if [ -x tools/tailwindcss ]; then \
	    ./tools/tailwindcss -i static/tailwind.in.css -o static/dist/tailwind.css --minify; \
	else \
	    echo "tools/tailwindcss not installed; download from https://github.com/tailwindlabs/tailwindcss/releases"; \
	fi

# ---------------------------------------------------------------------------
# Docker
# ---------------------------------------------------------------------------
.PHONY: docker-up
docker-up: ## Start the local Docker Compose stack.
	$(COMPOSE) -f $(COMPOSE_FILE) up -d

.PHONY: docker-down
docker-down: ## Stop the local Docker Compose stack.
	$(COMPOSE) -f $(COMPOSE_FILE) down

.PHONY: docker-logs
docker-logs: ## Tail Docker Compose logs.
	$(COMPOSE) -f $(COMPOSE_FILE) logs -f --tail=200

# ---------------------------------------------------------------------------
# License compliance
# ---------------------------------------------------------------------------
.PHONY: license-check
license-check: ## Fail the build if any dependency carries a disallowed license.
	bash scripts/license-check.sh

.PHONY: clean
clean: ## Remove build artefacts.
	rm -rf $(BIN_DIR) coverage.out coverage.html ext/dist

# ---------------------------------------------------------------------------
# WebExtension (Phase 21)
# ---------------------------------------------------------------------------
EXT_VERSION ?= 0.1.0

.PHONY: ext-package
ext-package: ## Build Joblantern WebExtension zips for Chrome and Firefox.
	VERSION=$(EXT_VERSION) bash ext/build.sh

.PHONY: ext-package-chrome
ext-package-chrome: ## Build only the Chrome zip.
	VERSION=$(EXT_VERSION) bash ext/build.sh chrome

.PHONY: ext-package-firefox
ext-package-firefox: ## Build only the Firefox zip.
	VERSION=$(EXT_VERSION) bash ext/build.sh firefox
