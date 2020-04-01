VERSION = $(shell git describe --dirty --tags --always)
DIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
PKGS = $(shell go list ./...)
TEST_PKGS = $(shell find . -type f -name "*_test.go" -printf '%h\n' | sort -u)
GOARGS = GOOS=linux GOARCH=amd64
GO_BUILD_ARGS = -buildmode=plugin -ldflags="-w -s"
GO_CONTAINER_BUILD_ARGS = -buildmode=plugin -ldflags="-w -s" -a -installsuffix cgo
GO_DEBUG_BUILD_ARGS = -buildmode=plugin -gcflags "all=-N -l"
PLUGIN_NAME = tls_interceptor.so
OUT_DIR = $(DIR)../../../plugins
DEBUG_PORT = 2345

.PHONY: deps format compile test cli-cover-report html-cover-report

all: compile test

deps:
	@go mod tidy
	@go build -buildmode=plugin -v $(DIR)...

format:
	@go fmt $(PKGS)

compile: deps
	@mkdir -p $(OUT_DIR)
ifdef DEBUG
	@echo 'Compiling for debugging...'
	@$(GOARGS) go build $(GO_DEBUG_BUILD_ARGS) -o $(OUT_DIR)/$(PLUGIN_NAME) $(DIR)
else ifdef CONTAINER
	@$(GOARGS) go build $(GO_CONTAINER_BUILD_ARGS) -o $(OUT_DIR)/$(PLUGIN_NAME) $(DIR)
else
	@$(GOARGS) go build $(GO_BUILD_ARGS) -o $(OUT_DIR)/$(PLUGIN_NAME) $(DIR)
endif

test:
	@go test -coverprofile=./cov-raw.out -v $(TEST_PKGS)
	@cat ./cov-raw.out | grep -v "generated" > ./cov.out

cli-cover-report:
	@go tool cover -func=cov.out

html-cover-report:
	@go tool cover -html=cov.out -o .coverage.html