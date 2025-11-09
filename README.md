# Containers API

A RESTful API for managing tasks and containers within user groups. Built with Go using clean architecture principles and CQRS pattern.

## Tech Stack

- **Go** 1.24.5
- **Framework:** Chi router
- **Database:** PostgreSQL (pgx/v5)
- **Authentication:** JWT
- **Logging:** Zerolog
- **Containerization:** Docker

## Features

- Task and container management
- User group organization with role-based permissions
- Command/Query separation (CQRS)
- Clean architecture with domain-driven design
- JWT authentication
- Comprehensive logging with request IDs

## Prerequisites
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- Make

## Quick Start

Start the application:
```sh
make start
```

Stop/remove containers:
```sh
make down
```

Rebuild:
```sh
make rebuild-docker
```

## Quick Commands for testing

| Command | Description |
|---------|-------------|
| `make test` | Run unit tests only (fast) |
| `make test-integration` | Run integration tests only |
| `make test-all` | Run all tests |
| `make test-coverage` | Generate coverage report |