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
	@echo "make test-db-start to start the test database"
	@echo "make test-db-stop to stop the test database"

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

# Testing targets
# - test: Unit tests only (internal/, pkg/, cmd/)
# - test-integration: Integration tests only (tests/integration/)
# - test-all: All tests (unit + integration)
# Note: All tests in tests/ folder are integration tests

test:
	@echo "Running unit tests (excluding tests/ folder)..."
	go test -v ./internal/... ./pkg/... ./cmd/...

test-db-start:
	@echo "Starting test database..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) -p $(CONTAINER_NAME) up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3

test-db-stop:
	@echo "Stopping test database..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) -p $(CONTAINER_NAME) stop postgres

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