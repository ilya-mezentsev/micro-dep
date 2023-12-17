ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
APPS_DIR := $(ROOT_DIR)/apps

USER_DIR := $(APPS_DIR)/user
STORE_DIR := $(APPS_DIR)/store
GOMODCACHE_DIR := $(APPS_DIR)/deps

TESTS_DIR := $(ROOT_DIR)/tests
TESTS_VENV := $(TESTS_DIR)/venv
TESTS_PIP := $(TESTS_VENV)/bin/pip
TESTS_PYTEST := $(TESTS_VENV)/bin/pytest

DOCKER_COMPOSE_ENTRYPOINT := $(ROOT_DIR)/docker-compose.yaml

ENTRYPOINT := cmd/main.go

mocks: store-mocks user-mocks

test: store-test user-test

tidy: store-tidy user-tidy

test-all: test e2e-test

e2e-test:
	$(TESTS_PYTEST) $(TESTS_DIR)

e2e-venv:
	python3 -m venv $(TESTS_VENV)

e2e-req:
	$(TESTS_PIP) install -r $(TESTS_DIR)/requirements.txt

e2e-setup: e2e-venv e2e-req

user-run:
	cd $(USER_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) CONFIG_PATH=./configs/main.json go run $(ENTRYPOINT)

user-tidy:
	cd $(USER_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) go mod tidy

user-mocks:
	cd $(USER_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) mockery

user-test:
	cd $(USER_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) go test -cover ./internal/... | { grep -v "no test files"; true; }

store-run:
	cd $(STORE_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) CONFIG_PATH=./configs/main.json go run $(ENTRYPOINT)

store-tidy:
	cd $(STORE_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) go mod tidy

store-mocks:
	cd $(STORE_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) mockery

store-test:
	cd $(STORE_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) go test -cover ./internal/... | { grep -v "no test files"; true; }

db-run:
	docker-compose -f $(DOCKER_COMPOSE_ENTRYPOINT) up -d db

down:
	docker-compose -f $(DOCKER_COMPOSE_ENTRYPOINT) down -v
