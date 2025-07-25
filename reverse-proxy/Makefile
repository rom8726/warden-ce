# Docker registry configuration
DOCKER_REGISTRY ?= docker.io

RED="\033[0;31m"
GREEN="\033[1;32m"
YELLOW="\033[0;33m"
NOCOLOR="\033[0m"

.DEFAULT_GOAL := help

.PHONY: help
help: ## Выводит это сообщение
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

.PHONY: docker-build
docker-build: ## Сборка Docker-образа
	@echo "\nBuilding Docker image..."
	@docker build \
		--build-arg TOOL_VERSION=${TOOL_VERSION} \
		--build-arg TOOL_BUILD_TIME=${TOOL_BUILD_TIME} \
		-t warden-reverse-proxy:latest -f ./Dockerfile .
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'warden-reverse-proxy' built successfully!"${NOCOLOR}

.PHONY: docker-push
docker-push: docker-build ## Тегирование и пуш Docker-образа
	@echo "\nTagging Docker image..."
	@docker tag warden-reverse-proxy:latest $(DOCKER_REGISTRY)/rom8726/warden-reverse-proxy:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/rom8726/warden-reverse-proxy:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker image with version $(TOOL_VERSION)..."; \
		docker tag warden-reverse-proxy:latest $(DOCKER_REGISTRY)/rom8726/warden-reverse-proxy:$(TOOL_VERSION); \
		                if [ $$? -ne 0 ] ; then \
                        echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
                        exit 1 ; \
                fi; \
                echo "\nPushing Docker image with version $(TOOL_VERSION) to Docker Registry..."; \
                docker push $(DOCKER_REGISTRY)/rom8726/warden-reverse-proxy:$(TOOL_VERSION); \
                if [ $$? -ne 0 ] ; then \
                        echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker image pushed to registry successfully!"${NOCOLOR}
