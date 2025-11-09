DEV_ENV_SETUP_FOLDER ?= ./dev-env
DOCKER_COMPOSE_FILE ?= $(DEV_ENV_SETUP_FOLDER)/docker-compose.yml
CONTAINER_NAME ?= "task-containers"
VERSION ?= $(shell git rev-parse --short HEAD)

help:
	@echo "make version to get the current version"
	@echo "make start to start go-api server"
	@echo "make build"
	@echo "make rebuild-docker"
	@echo "make logs"
	@echo "make down to remove docker containers"
	@echo "make clean to remove containers, images, and volumes completely"
	@echo "make test to run unit tests only"
	@echo "make test-integration to run integration tests"
	@echo "make test-all to run all tests (unit + integration)"

version:
	@echo $(VERSION)

start:
	@echo "Starting app..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) -p $(CONTAINER_NAME) up --build -d

down:
	@echo "Stopping app..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) -p $(CONTAINER_NAME) down
	@rm -rf $(DEV_ENV_SETUP_FOLDER)/postgres-data
	@echo "Docker containers removed and postgres data cleaned up."

force-clean:
	@echo "Removing all containers, images, and data..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) -p $(CONTAINER_NAME) down --rmi all --volumes --remove-orphans
	@rm -rf $(DEV_ENV_SETUP_FOLDER)/postgres-data
	@echo "All Docker containers, images, volumes removed and postgres data cleaned up."
build:
	go build -v ./...

rebuild-docker:
	@echo "Rebuilding Docker containers..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) -p $(CONTAINER_NAME) down
	@docker compose -f $(DOCKER_COMPOSE_FILE) -p $(CONTAINER_NAME) build --no-cache
	@docker compose -f $(DOCKER_COMPOSE_FILE) -p $(CONTAINER_NAME) up -d
	@echo "Docker containers rebuilt!"

watch: start
	@echo "Watching for file changes..."
	@docker-compose watch

logs:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) -p $(CONTAINER_NAME) logs -f

test:
	@echo "Running unit tests (excluding tests/ folder)..."
	go test -v ./internal/... ./pkg/... ./cmd/...

test-integration: test-db-start
	@echo "Running integration tests..."
	go test -v ./tests/integration/... -count=1
	@echo "Integration tests completed!"

test-all: test-db-start
	@echo "Running all tests (unit + integration)..."
	go test -v ./... -count=1
	@echo "All tests completed!"

test-coverage: test-db-start
	@echo "Running tests with coverage..."
	go test -v ./internal/... ./pkg/... ./cmd/... ./tests/... -coverprofile=coverage.out -count=1
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"