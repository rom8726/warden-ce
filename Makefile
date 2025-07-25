TOOLS_DIR=dev/tools
NAMESPACE=warden
GOCI_VERSION=v1.64.8
MOCKERY_VERSION=v2.53.3
TOOLS_DIR_ABS=${PWD}/${TOOLS_DIR}
BIN_OUTPUT_DIR=bin
GOLANGCI_LINT=${TOOLS_DIR}/golangci-lint
MOCKERY=${TOOLS_DIR}/mockery
TPARSE=${TOOLS_DIR}/tparse
GOCMD=go
GOBUILD=$(GOCMD) build
GOPROXY=https://proxy.golang.org,direct
TOOL_VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
TOOL_BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
OS=$(shell uname -s)
IGNORE_TEST_DIRS=cmd,ui,test_mocks,specs,scripts,migrations,docs,dev,internal/generated,tests/runner,internal/contract,pkg/metrics

# Docker registry configuration
DOCKER_REGISTRY ?= docker.io

LD_FLAGS="-w -s -X 'github.com/rom8726/warden/internal/version.Version=${TOOL_VERSION}' -X 'github.com/rom8726/warden/internal/version.BuildTime=${TOOL_BUILD_TIME}' -X 'github.com/rom8726/warden/internal/installer.DockerRegistry=${DOCKER_REGISTRY}'"

RED="\033[0;31m"
GREEN="\033[1;32m"
YELLOW="\033[0;33m"
NOCOLOR="\033[0m"

.DEFAULT_GOAL := help

#
# Extra targets
#
-include dev/dev.mk
-include docker.mk
-include profile.mk

#
# Local targets
#

.PHONY: help
help: ## Print this message
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

.PHONY: .install-linter
.install-linter:
	@[ -f $(GOLANGCI_LINT) ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_DIR) $(GOCI_VERSION)

.PHONY: .install-mockery
.install-mockery:
	@[ -f $(MOCKERY) ] || GOBIN=$(TOOLS_DIR_ABS) go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)

.PHONY: .install_tparse
.install_tparse:
	@[ -f $(TPARSE) ] || GOBIN=$(TOOLS_DIR_ABS) go install github.com/mfridman/tparse@latest

.PHONY: setup
setup: .install-linter .install-mockery ## Setup development environment
	@echo "\nCreate .env files in dev/ directory"
	@cp dev/config.env.example dev/config.env
	@cp dev/compose.env.example dev/compose.env

	@echo
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"OK"${NOCOLOR}

.PHONY: lint
lint: .install-linter ## Run linter
	@$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY: test
test: .install_tparse ## Run unit tests
	@IGNORE_PATTERN="$$(echo $(IGNORE_TEST_DIRS) | tr ',' '|')" ; \
	PKGS="$$(go list ./... | grep -Ev "/($$IGNORE_PATTERN)(/|$$)")" ; \
	go test -json -cover -coverprofile=coverage.out -v $$PKGS > test.out || true
	@go tool cover -html=coverage.out -o coverage.html
	@$(TPARSE) -all -file=test.out

.PHONY: test.integration
test.integration: .install_tparse ## Run unit and integration tests
	@IGNORE_PATTERN="$$(echo $(IGNORE_TEST_DIRS) | tr ',' '|')" ; \
	PKGS="$$(go list -tags=integration ./... | grep -Ev "/($$IGNORE_PATTERN)(/|$$)")" ; \
	go test -tags=integration -json -cover -coverprofile=coverage.out -v $$PKGS > test.out || true
	@go tool cover -html=coverage.out -o coverage.html
	@$(TPARSE) -all -file=test.out

.PHONY: build-backend
build-backend: ## Build backend binary
	@echo "\nBuilding backend binary..."
	@echo
	go env -w GOPROXY=${GOPROXY}
	go env -w GOPRIVATE=${GOPRIVATE}

	CGO_ENABLED=0 $(GOBUILD) -trimpath -ldflags=$(LD_FLAGS) -o ${BIN_OUTPUT_DIR}/app ./cmd/backend

.PHONY: build-ingest-server
build-ingest-server: ## Build ingest server binary
	@echo "\nBuilding ingest server binary..."
	@echo
	go env -w GOPROXY=${GOPROXY}
	go env -w GOPRIVATE=${GOPRIVATE}

	CGO_ENABLED=0 $(GOBUILD) -trimpath -ldflags=$(LD_FLAGS) -o ${BIN_OUTPUT_DIR}/app ./cmd/ingest-server

.PHONY: build-envelope-consumer
build-envelope-consumer: ## Build envelope consumer binary
	@echo "\nBuilding envelope consumer binary..."
	@echo
	go env -w GOPROXY=${GOPROXY}
	go env -w GOPRIVATE=${GOPRIVATE}

	CGO_ENABLED=0 $(GOBUILD) -trimpath -ldflags=$(LD_FLAGS) -o ${BIN_OUTPUT_DIR}/app ./cmd/envelope-consumer

.PHONY: build-issue-notificator
build-issue-notificator: ## Build issue notificator binary
	@echo "\nBuilding issue notificator binary..."
	@echo
	go env -w GOPROXY=${GOPROXY}
	go env -w GOPRIVATE=${GOPRIVATE}

	CGO_ENABLED=0 $(GOBUILD) -trimpath -ldflags=$(LD_FLAGS) -o ${BIN_OUTPUT_DIR}/app ./cmd/issue-notificator

.PHONY: build-user-notificator
build-user-notificator: ## Build user notificator binary
	@echo "\nBuilding user notificator binary..."
	@echo
	go env -w GOPROXY=${GOPROXY}
	go env -w GOPRIVATE=${GOPRIVATE}

	CGO_ENABLED=0 $(GOBUILD) -trimpath -ldflags=$(LD_FLAGS) -o ${BIN_OUTPUT_DIR}/app ./cmd/user-notificator

.PHONY: build-scheduler
build-scheduler: ## Build scheduler binary
	@echo "\nBuilding scheduler binary..."
	@echo
	go env -w GOPROXY=${GOPROXY}
	go env -w GOPRIVATE=${GOPRIVATE}

	CGO_ENABLED=0 $(GOBUILD) -trimpath -ldflags=$(LD_FLAGS) -o ${BIN_OUTPUT_DIR}/app ./cmd/scheduler

.PHONY: build-frontend
build-frontend: ## Build frontend with version
	@echo "\nBuilding frontend with version..."
	@cd ui && make build
	@if [ $$? -ne 0 ] ; then \
		echo -e ${RED}"Frontend build FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Frontend built successfully!"${NOCOLOR}

.PHONY: build-installer
build-installer: ## Build installer binary
	@echo "\nBuilding installer binary..."
	@echo
	go env -w GOPROXY=${GOPROXY}
	go env -w GOPRIVATE=${GOPRIVATE}

	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} $(GOBUILD) -trimpath -ldflags=$(LD_FLAGS) -o ${BIN_OUTPUT_DIR}/installer ./cmd/installer

.PHONY: build-demo-error-sender
build-demo-error-sender: ## Build demo error sender binary for Linux AMD64
	@echo "\nBuilding demo error sender binary for Linux AMD64..."
	@echo
	go env -w GOPROXY=${GOPROXY}
	go env -w GOPRIVATE=${GOPRIVATE}

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath -o ${BIN_OUTPUT_DIR}/demo_error_sender ./demo/demo_error_sender.go

.PHONY: mocks
mocks: .install-mockery ## Generate mocks with mockery
	./dev/tools/mockery

.PHONY: format
format: ## Format go code
	go fmt ./...

.PHONY: generate-backend
generate-backend: ## Generate backend by OpenAPI specification
	@docker run --rm \
      --volume ".:/workspace" \
      ghcr.io/ogen-go/ogen:latest --target workspace/internal/generated/server --clean workspace/specs/server.yml

.PHONY: generate-ingest-server
generate-ingest-server: ## Generate ingest server by OpenAPI specification
	@docker run --rm \
      --volume ".:/workspace" \
      ghcr.io/ogen-go/ogen:latest --target workspace/internal/generated/ingestserver --clean workspace/specs/ingest_server.yml
