ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
APPS_DIR := $(ROOT_DIR)/apps

USER_DIR := $(APPS_DIR)/user
STORE_DIR := $(APPS_DIR)/store
DIAGRAM_DIR := $(APPS_DIR)/diagram
GOMODCACHE_DIR := $(APPS_DIR)/deps

TESTS_DIR := $(ROOT_DIR)/tests
TESTS_VENV := $(TESTS_DIR)/venv
TESTS_PIP := $(TESTS_VENV)/bin/pip
TESTS_PYTEST := $(TESTS_VENV)/bin/pytest

DOCKER_COMPOSE_ENTRYPOINT := $(ROOT_DIR)/docker-compose.yaml

ENTRYPOINT := cmd/main.go
COMPILATION_OUTPUT := compiled/main

mocks: store-mocks user-mocks diagram-mocks

build: store-build user-build diagram-build

run: down build
	docker-compose -f $(DOCKER_COMPOSE_ENTRYPOINT) up

down:
	docker-compose -f $(DOCKER_COMPOSE_ENTRYPOINT) down -v

test: store-test user-test diagram-test

tidy: store-tidy user-tidy diagram-tidy

test-all: test e2e-test

e2e-test:
	$(TESTS_PYTEST) $(TESTS_DIR)

e2e-venv:
	python3 -m venv $(TESTS_VENV)

e2e-req:
	$(TESTS_PIP) install -r $(TESTS_DIR)/requirements.txt

e2e-setup: e2e-venv e2e-req

user-run:
	@$(MAKE) --no-print-directory APP_DIR=$(USER_DIR) app-run

user-build:
	@$(MAKE) --no-print-directory APP_DIR=$(USER_DIR) app-build

user-tidy:
	@$(MAKE) --no-print-directory APP_DIR=$(USER_DIR) app-tidy

user-mocks:
	@$(MAKE) --no-print-directory APP_DIR=$(USER_DIR) app-mocks

user-test:
	@$(MAKE) --no-print-directory APP_DIR=$(USER_DIR) app-test

store-run:
	@$(MAKE) --no-print-directory APP_DIR=$(STORE_DIR) app-run

store-build:
	@$(MAKE) --no-print-directory APP_DIR=$(STORE_DIR) app-build

store-tidy:
	@$(MAKE) --no-print-directory APP_DIR=$(STORE_DIR) app-tidy

store-mocks:
	@$(MAKE) --no-print-directory APP_DIR=$(STORE_DIR) app-mocks

store-test:
	@$(MAKE) --no-print-directory APP_DIR=$(STORE_DIR) app-test

diagram-run:
	@$(MAKE) --no-print-directory APP_DIR=$(DIAGRAM_DIR) app-run

diagram-build:
	@$(MAKE) --no-print-directory APP_DIR=$(DIAGRAM_DIR) app-build

diagram-tidy:
	@$(MAKE) --no-print-directory APP_DIR=$(DIAGRAM_DIR) app-tidy

diagram-mocks:
	@$(MAKE) --no-print-directory APP_DIR=$(DIAGRAM_DIR) app-mocks

diagram-test:
	@$(MAKE) --no-print-directory APP_DIR=$(DIAGRAM_DIR) app-test

app-run:
	cd $(APP_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) CONFIG_PATH=./configs/main.json go run $(ENTRYPOINT)

app-build:
	cd $(APP_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) go build -o $(COMPILATION_OUTPUT) $(ENTRYPOINT)

app-tidy:
	cd $(APP_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) go mod tidy

app-mocks:
	cd $(APP_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) mockery

app-test:
	cd $(APP_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) go test -cover ./internal/... | { grep -v "no test files"; true; }

db-run:
	docker-compose -f $(DOCKER_COMPOSE_ENTRYPOINT) up -d db
