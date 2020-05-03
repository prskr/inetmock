VERSION = $(shell git describe --dirty --tags --always)
DIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
BUILD_PATH = $(DIR)/main.go
PKGS = $(shell go list ./...)
TEST_PKGS = $(shell go list ./...)
PROTO_FILES = $(shell find $(DIR)api/ -type f -name "*.proto")
GOARGS = GOOS=linux GOARCH=amd64
GO_BUILD_ARGS = -ldflags="-w -s"
GO_CONTAINER_BUILD_ARGS = -ldflags="-w -s" -a -installsuffix cgo
GO_DEBUG_BUILD_ARGS = -gcflags "all=-N -l"
BINARY_NAME = inetmock
PLUGINS = $(wildcard $(DIR)plugins/*/.)
DEBUG_PORT = 2345
DEBUG_ARGS?= --development-logs=true
CONTAINER_BUILDER ?= podman
DOCKER_IMAGE ?= inetmock

.PHONY: clean all format deps update-deps compile debug generate protoc snapshot-release test cli-cover-report html-cover-report plugins $(PLUGINS) $(GO_GEN_FILES)
all: clean format compile test plugins

clean:
	@find $(DIR) -type f \( -name "*.out" -or -name "*.so" \) -exec rm -f {} \;
	@rm -rf $(DIR)*.so
	@rm -f $(DIR)$(BINARY_NAME) $(DIR)main

format:
	@go fmt $(PKGS)

deps:
	@go build -v $(BUILD_PATH)

update-deps:
	@go mod tidy
	@go get -u

compile: deps
ifdef DEBUG
	@echo 'Compiling for debugging...'
	@$(GOARGS) go build $(GO_DEBUG_BUILD_ARGS) -o $(DIR)$(BINARY_NAME) $(BUILD_PATH)
else ifdef CONTAINER
	@echo 'Compiling for container usage...'
	@$(GOARGS) go build $(GO_CONTAINER_BUILD_ARGS) -o $(DIR)$(BINARY_NAME) $(BUILD_PATH)
else
	@echo 'Compiling for normal Linux env...'
	@$(GOARGS) go build $(GO_BUILD_ARGS) -o $(DIR)$(BINARY_NAME) $(BUILD_PATH)
endif

debug: export INETMOCK_PLUGINS_DIRECTORY = $(DIR)
debug:
	dlv debug $(DIR) \
		--headless \
		--listen=:2345 \
		--api-version=2 \
		-- $(DEBUG_ARGS)

generate:
	@go generate ./...
	@protoc --proto_path $(DIR)api/ --go_out=plugins=grpc:internal/rpc --go_opt=paths=source_relative $(shell find $(DIR)api/ -type f -name "*.proto")

snapshot-release:
	@goreleaser release --snapshot --skip-publish --rm-dist

container:
	@$(CONTAINER_BUILDER) build -t $(DOCKER_IMAGE):latest -f $(DIR)Dockerfile $(DIR)

test:
	@go test -coverprofile=./cov-raw.out -v $(TEST_PKGS)
	@cat ./cov-raw.out | grep -v "generated" > ./cov.out

cli-cover-report: test
	@go tool cover -func=cov.out

html-cover-report: test
	@go tool cover -html=cov.out -o .coverage.html

plugins: $(PLUGINS)
$(PLUGINS):
	$(MAKE) -C $@
