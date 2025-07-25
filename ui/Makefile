# Docker registry configuration
DOCKER_REGISTRY ?= chudov.cr.cloud.ru

RED="\033[0;31m"
GREEN="\033[1;32m"
YELLOW="\033[0;33m"
NOCOLOR="\033[0m"

# Version variables - same logic as main Makefile
FRONTEND_VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
FRONTEND_BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

.DEFAULT_GOAL := help

.PHONY: help
help: ## Выводит это сообщение
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

.PHONY: lint
lint: ## Запуск линтера
	@npm run lint

.PHONY: test
test: ## Запуск тестов
	@npm run test

.PHONY: test-coverage
test-coverage: ## Запуск тестов с отчетом о покрытии
	@echo "Running tests and generating coverage report..."
	@npm run test

.PHONY: build
build: ## Сборка фронтенда с версией из git
	@echo "Building frontend with version: ${FRONTEND_VERSION}"
	@echo "Build time: ${FRONTEND_BUILD_TIME}"
	@VITE_VERSION=${FRONTEND_VERSION} VITE_BUILD_TIME=${FRONTEND_BUILD_TIME} npm run build
	@if [ $$? -ne 0 ] ; then \
		echo -e ${RED}"Build FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Frontend built successfully!"${NOCOLOR}

PHONY: generate-api-client
generate-api-client: ## Сгенерировать API клиент по спецификации
	cp ../specs/server.yml specs/server.yml
	@docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:latest generate -i /local/specs/server.yml -g typescript-axios -o /local/src/generated/api/client

	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"\nOK"${NOCOLOR}

.PHONY: docker-build
docker-build: ## Сборка универсального Docker-образа
	@echo "\nBuilding universal Docker image..."
	@docker build \
		--build-arg VITE_VERSION=$(FRONTEND_VERSION) \
		--build-arg VITE_BUILD_TIME=$(FRONTEND_BUILD_TIME) \
		-t warden-frontend:latest \
		-f ./Dockerfile .
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Universal Docker image 'warden-frontend' built successfully!"${NOCOLOR}
	@echo ${YELLOW}"Note: API URL will be configured at runtime via WARDEN_API_BASE_URL environment variable"${NOCOLOR}

.PHONY: docker-push
docker-push: docker-build ## Тегирование и пуш Docker-образа
	@echo "\nTagging Docker image..."
	@docker tag warden-frontend:latest $(DOCKER_REGISTRY)/warden-frontend:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/warden-frontend:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(VITE_VERSION)" ] && [ "$(VITE_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker image with version $(VITE_VERSION)..."; \
		docker tag warden-frontend:latest $(DOCKER_REGISTRY)/warden-frontend:$(VITE_VERSION); \
		                if [ $$? -ne 0 ] ; then \
                        echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
                        exit 1 ; \
                fi; \
                echo "\nPushing Docker image with version $(VITE_VERSION) to Docker Registry..."; \
                docker push $(DOCKER_REGISTRY)/warden-frontend:$(VITE_VERSION); \
                if [ $$? -ne 0 ] ; then \
                        echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker image pushed to registry successfully!"${NOCOLOR}
