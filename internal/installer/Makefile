_COMPOSE=docker compose -f docker-compose.yml --project-name warden --env-file platform.env

.DEFAULT_GOAL := help

.PHONY: help
help: ## Print this message
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

.PHONY: up
up: ## Up the environment in docker compose
	${_COMPOSE} up -d

.PHONY: down
down: ## Down the environment in docker compose
	${_COMPOSE} down --remove-orphans

.PHONY: pull
pull: ## Pull images from remote Docker registry
	${_COMPOSE} pull
