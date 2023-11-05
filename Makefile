ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
APPS_DIR := $(ROOT_DIR)/apps

USER_DIR := $(APPS_DIR)/user
STORE_DIR := $(APPS_DIR)/store
GOMODCACHE_DIR := $(APPS_DIR)/deps

ENTRYPOINT := cmd/main.go

mocks: store-mocks

test: store-test

user-run:
	cd $(USER_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) go run $(ENTRYPOINT)

store-run:
	cd $(STORE_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) go run $(ENTRYPOINT)

store-tidy:
	cd $(STORE_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) go mod tidy

store-mocks:
	cd $(STORE_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) mockery

store-test:
	cd $(STORE_DIR) && GOMODCACHE=$(GOMODCACHE_DIR) go test -cover ./internal/... | { grep -v "no test files"; true; }
