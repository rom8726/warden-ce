RED="\033[0;31m"
GREEN="\033[1;32m"
NOCOLOR="\033[0m"

.PHONY: docker-build-backend
docker-build-backend: ## Building Docker image for backend (scratch + curl)
	@echo "\nBuilding Docker image (scratch + curl)..."
	@docker build \
		--build-arg TOOL_VERSION=${TOOL_VERSION} \
		--build-arg TOOL_BUILD_TIME=${TOOL_BUILD_TIME} \
		--build-arg DOCKER_REGISTRY=${DOCKER_REGISTRY} \
		-t warden-backend:latest -f ./Dockerfile.backend .
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'warden-backend' built successfully!"${NOCOLOR}

.PHONY: docker-push-backend
docker-push-backend: docker-build-backend ## Tagging and pushing Docker image for backend
	@echo "\nTagging Docker backend image..."
	@docker tag warden-backend:latest $(DOCKER_REGISTRY)/warden-backend:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker backend image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/warden-backend:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker backend image with version $(TOOL_VERSION)..."; \
		docker tag warden-backend:latest $(DOCKER_REGISTRY)/warden-backend:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker backend image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/warden-backend:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker backend image pushed to registry successfully!"${NOCOLOR}

.PHONY: docker-build-ingest-server
docker-build-ingest-server: ## Building Docker image for ingest-server (scratch + curl)
	@echo "\nBuilding Docker image (scratch + curl)..."
	@docker build \
		--build-arg TOOL_VERSION=${TOOL_VERSION} \
		--build-arg TOOL_BUILD_TIME=${TOOL_BUILD_TIME} \
		--build-arg DOCKER_REGISTRY=${DOCKER_REGISTRY} \
		-t warden-ingest-server:latest -f ./Dockerfile.ingest_server .
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'warden-ingest-server' built successfully!"${NOCOLOR}

.PHONY: docker-push-ingest-server
docker-push-ingest-server: docker-build-ingest-server ## Tagging and pushing Docker image for ingest-server
	@echo "\nTagging Docker ingest-server image..."
	@docker tag warden-ingest-server:latest $(DOCKER_REGISTRY)/warden-ingest-server:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker ingest-server image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/warden-ingest-server:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker ingest-server image with version $(TOOL_VERSION)..."; \
		docker tag warden-ingest-server:latest $(DOCKER_REGISTRY)/warden-ingest-server:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker ingest-server image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/warden-ingest-server:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker ingest-server image pushed to registry successfully!"${NOCOLOR}

.PHONY: docker-build-envelope-consumer
docker-build-envelope-consumer: ## Building Docker image for envelope-consumer (scratch + curl)
	@echo "\nBuilding Docker image (scratch + curl)..."
	@docker build \
		--build-arg TOOL_VERSION=${TOOL_VERSION} \
		--build-arg TOOL_BUILD_TIME=${TOOL_BUILD_TIME} \
		--build-arg DOCKER_REGISTRY=${DOCKER_REGISTRY} \
		-t warden-envelope-consumer:latest -f ./Dockerfile.envelope_consumer .
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'warden-envelope-consumer' built successfully!"${NOCOLOR}

.PHONY: docker-push-envelope-consumer
docker-push-envelope-consumer: docker-build-envelope-consumer ## Tagging and pushing Docker image for envelope-consumer
	@echo "\nTagging Docker envelope-consumer image..."
	@docker tag warden-envelope-consumer:latest $(DOCKER_REGISTRY)/warden-envelope-consumer:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker envelope-consumer image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/warden-envelope-consumer:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker envelope-consumer image with version $(TOOL_VERSION)..."; \
		docker tag warden-envelope-consumer:latest $(DOCKER_REGISTRY)/warden-envelope-consumer:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker envelope-consumer image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/warden-envelope-consumer:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker envelope-consumer image pushed to registry successfully!"${NOCOLOR}

.PHONY: docker-build-issue-notificator
docker-build-issue-notificator: ## Building Docker image for issue-notificator (scratch + curl)
	@echo "\nBuilding Docker image (scratch + curl)..."
	@docker build \
		--build-arg TOOL_VERSION=${TOOL_VERSION} \
		--build-arg TOOL_BUILD_TIME=${TOOL_BUILD_TIME} \
		--build-arg DOCKER_REGISTRY=${DOCKER_REGISTRY} \
		-t warden-issue-notificator:latest -f ./Dockerfile.issue_notificator .
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'warden-issue-notificator' built successfully!"${NOCOLOR}

.PHONY: docker-push-issue-notificator
docker-push-issue-notificator: docker-build-issue-notificator ## Tagging and pushing Docker image for issue-notificator
	@echo "\nTagging Docker issue-notificator image..."
	@docker tag warden-issue-notificator:latest $(DOCKER_REGISTRY)/warden-issue-notificator:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker issue-notificator image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/warden-issue-notificator:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker issue-notificator image with version $(TOOL_VERSION)..."; \
		docker tag warden-issue-notificator:latest $(DOCKER_REGISTRY)/warden-issue-notificator:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker issue-notificator image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/warden-issue-notificator:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker issue-notificator image pushed to registry successfully!"${NOCOLOR}

.PHONY: docker-build-user-notificator
docker-build-user-notificator: ## Building Docker image for user-notificator (scratch + curl)
	@echo "\nBuilding Docker image (scratch + curl)..."
	@docker build \
		--build-arg TOOL_VERSION=${TOOL_VERSION} \
		--build-arg TOOL_BUILD_TIME=${TOOL_BUILD_TIME} \
		--build-arg DOCKER_REGISTRY=${DOCKER_REGISTRY} \
		-t warden-user-notificator:latest -f ./Dockerfile.user_notificator .
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'warden-user-notificator' built successfully!"${NOCOLOR}

.PHONY: docker-push-user-notificator
docker-push-user-notificator: docker-build-user-notificator ## Tagging and pushing Docker image for user-notificator
	@echo "\nTagging Docker user-notificator image..."
	@docker tag warden-user-notificator:latest $(DOCKER_REGISTRY)/warden-user-notificator:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker user-notificator image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/warden-user-notificator:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker user-notificator image with version $(TOOL_VERSION)..."; \
		docker tag warden-user-notificator:latest $(DOCKER_REGISTRY)/warden-user-notificator:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker user-notificator image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/warden-user-notificator:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker user-notificator image pushed to registry successfully!"${NOCOLOR}

.PHONY: docker-build-scheduler
docker-build-scheduler: ## Building Docker image for scheduler (scratch + curl)
	@echo "\nBuilding Docker image (scratch + curl)..."
	@docker build \
		--build-arg TOOL_VERSION=${TOOL_VERSION} \
		--build-arg TOOL_BUILD_TIME=${TOOL_BUILD_TIME} \
		--build-arg DOCKER_REGISTRY=${DOCKER_REGISTRY} \
		-t warden-scheduler:latest -f ./Dockerfile.scheduler .
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'warden-scheduler' built successfully!"${NOCOLOR}

.PHONY: docker-push-scheduler
docker-push-scheduler: docker-build-scheduler ## Tagging and pushing Docker image for scheduler
	@echo "\nTagging Docker scheduler image..."
	@docker tag warden-scheduler:latest $(DOCKER_REGISTRY)/warden-scheduler:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker scheduler image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/warden-scheduler:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker scheduler image with version $(TOOL_VERSION)..."; \
		docker tag warden-scheduler:latest $(DOCKER_REGISTRY)/warden-scheduler:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker scheduler image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/warden-scheduler:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker scheduler image pushed to registry successfully!"${NOCOLOR}

.PHONY: docker-build-pgbouncer
docker-build-pgbouncer: ## Building Docker image for pgbouncer
	@echo "\nBuilding Docker image pgbouncer..."
	@docker build \
		--build-arg TOOL_VERSION=${TOOL_VERSION} \
		--build-arg TOOL_BUILD_TIME=${TOOL_BUILD_TIME} \
		-t warden-pgbouncer:latest -f ./pgbouncer/Dockerfile ./pgbouncer
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'warden-pgbouncer' built successfully!"${NOCOLOR}

.PHONY: docker-push-pgbouncer
docker-push-pgbouncer: docker-build-pgbouncer ## Tagging and pushing Docker image for pgbouncer
	@echo "\nTagging Docker pgbouncer image..."
	@docker tag warden-pgbouncer:latest $(DOCKER_REGISTRY)/warden-pgbouncer:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi

	@echo "\nPushing Docker pgbouncer image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/warden-pgbouncer:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker pgbouncer image with version $(TOOL_VERSION)..."; \
		docker tag warden-pgbouncer:latest $(DOCKER_REGISTRY)/warden-pgbouncer:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker pgbouncer image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/warden-pgbouncer:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker pgbouncer image pushed to registry successfully!"${NOCOLOR}

.PHONY: docker-build-frontend
docker-build-frontend: ## Building Docker image for frontend
	@echo "\nBuilding Docker image frontend..."
	@cd ui && TOOL_VERSION=${TOOL_VERSION} TOOL_BUILD_TIME=${TOOL_BUILD_TIME} VITE_IS_DEMO=${VITE_IS_DEMO} make docker-build
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'warden-frontend' built successfully!"${NOCOLOR}

.PHONY: docker-build-frontend-demo
docker-build-frontend-demo: ## Building Docker image frontend for demo
	@echo "\nBuilding Docker image frontend..."
	@cd ui && TOOL_VERSION=${TOOL_VERSION} TOOL_BUILD_TIME=${TOOL_BUILD_TIME} VITE_IS_DEMO=true make docker-build
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'warden-frontend' (demo) built successfully!"${NOCOLOR}

.PHONY: docker-push-frontend
docker-push-frontend: docker-build-frontend ## Tagging and pushing Docker image for frontend
	@echo "\nTagging Docker frontend image..."
	@docker tag warden-frontend:latest $(DOCKER_REGISTRY)/warden-frontend:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker frontend image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/warden-frontend:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker frontend image with version $(TOOL_VERSION)..."; \
		docker tag warden-frontend:latest $(DOCKER_REGISTRY)/warden-frontend:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker frontend image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/warden-frontend:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker frontend image pushed to registry successfully!"${NOCOLOR}

.PHONY: docker-build-reverse-proxy
docker-build-reverse-proxy: ## Building Docker image for reverse-proxy
	@echo "\nBuilding Docker image reverse-proxy..."
	@cd reverse-proxy && TOOL_VERSION=${TOOL_VERSION} TOOL_BUILD_TIME=${TOOL_BUILD_TIME} make docker-build
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'warden-reverse-proxy' built successfully!"${NOCOLOR}

.PHONY: docker-push-reverse-proxy
docker-push-reverse-proxy: docker-build-reverse-proxy ## Tagging and pushing Docker image for reverse-proxy
	@echo "\nTagging Docker reverse-proxy image..."
	@docker tag warden-reverse-proxy:latest $(DOCKER_REGISTRY)/warden-reverse-proxy:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker reverse-proxy image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/warden-reverse-proxy:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker reverse-proxy image with version $(TOOL_VERSION)..."; \
		docker tag warden-reverse-proxy:latest $(DOCKER_REGISTRY)/warden-reverse-proxy:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker reverse-proxy image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/warden-reverse-proxy:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker reverse-proxy image pushed to registry successfully!"${NOCOLOR}
