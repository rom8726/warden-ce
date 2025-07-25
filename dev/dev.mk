_COMPOSE=docker compose -f dev/docker-compose.yml --project-name ${NAMESPACE} --env-file dev/platform.env

dev-up: ## Up the environment in docker compose
	${_COMPOSE} up -d

dev-down: ## Down the environment in docker compose
	${_COMPOSE} down --remove-orphans

dev-clean: ## Down the environment in docker compose with image cleanup
	${_COMPOSE} down --remove-orphans -v --rmi all

dev-logs: ## Show logs from all containers
	${_COMPOSE} logs -f

dev-email-test: ## Send a test email through MailHog
	@echo "Sending test email..."
	@sh dev/send-test-email.sh
	@echo "Test email sent. Check MailHog UI at http://localhost:8025"

dev-cert: ## Generates certificate for domain
	@mkdir -p dev/nginx/ssl
	@openssl req -newkey rsa:4096 -keyout dev/nginx/ssl/localhost.key -out dev/nginx/ssl/localhost.csr -nodes -subj "/C=RU/ST=Moscow/L=Moscow/O=Warden/OU=Warden/CN=warden"
	@openssl x509 -req -in dev/nginx/ssl/localhost.csr -signkey dev/nginx/ssl/localhost.key -out dev/nginx/ssl/localhost.crt -days 365

dev-build-proxy: ## Building warden-reverse-proxy
	${_COMPOSE} build warden-reverse-proxy

dev-build-frontend: ## Building warden-frontend
	${_COMPOSE} build warden-frontend

dev-build-backend: ## Building warden-backend
	${_COMPOSE} build warden-backend

dev-build-envelope-consumer: ## Building warden-envelope-consumer
	${_COMPOSE} build warden-envelope-consumer

dev-build-all: ## Building all components
	${_COMPOSE} build
