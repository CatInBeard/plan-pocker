# Makefile for managing Docker Compose local project

COMPOSE_FILE = docker-compose.yaml
ENV_FILE = local/local.env

.PHONY: help up down rebuild migrate fresh-migrate restart-service clear-local

help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

up: # Start the services defined in docker-compose.yml
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d

down: # Stop the services defined in docker-compose.yml
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) down

rebuild: # Rebuild the services and recreate the database
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) down
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up --build -d

migrate: # Run pending migrations
	docker compose --env-file $(ENV_FILE) up flyway 

clear-local: # Delete local data for all containers (db, cache, app logs, nginx logs, redisinsights)
	rm ./local/data