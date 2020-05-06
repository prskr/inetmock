VERSION = $(shell git describe --dirty --tags --always)
DIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
SERVER_BUILD_PATH = github.com/baez90/inetmock/cmd/inetmock
CLI_BUILD_PATH = github.com/baez90/inetmock/cmd/imctl
PKGS = $(shell go list ./...)
TEST_PKGS = $(shell go list ./...)
PROTO_FILES = $(shell find $(DIR)api/ -type f -name "*.proto")
GOARGS = GOOS=linux GOARCH=amd64
GO_BUILD_ARGS = -ldflags='-w -s'
GO_CONTAINER_BUILD_ARGS = -ldflags='-w -s' -a -installsuffix cgo
GO_DEBUG_BUILD_ARGS = -gcflags "all=-N -l"
SERVER_BINARY_NAME = inetmock
CLI_BINARY_NAME = imctl
PLUGINS = $(wildcard $(DIR)plugins/*/.)
DEBUG_PORT = 2345
DEBUG_ARGS?= --development-logs=true
CONTAINER_BUILDER ?= podman
DOCKER_IMAGE ?= inetmock

.PHONY: clean all format deps update-deps compile compile-server compile-cli debug generate protoc snapshot-release test cli-cover-report html-cover-report plugins $(PLUGINS) $(GO_GEN_FILES)
all: clean format generate compile test plugins

clean:
	@find $(DIR) -type f \( -name "*.out" -or -name "*.so" \) -exec rm -f {} \;
	@rm -rf $(DIR)*.so
	@find $(DIR) -type f -name "*.pb.go" -exec rm -f {} \;
	@find $(DIR) -type f -name "*.mock.go" -exec rm -f {} \;
	@rm -f $(DIR)$(SERVER_BINARY_NAME) $(DIR)$(CLI_BINARY_NAME) $(DIR)main

format:
	@go fmt $(PKGS)

deps:
	@go build -v $(SERVER_BUILD_PATH)

update-deps:
	@go mod tidy
	@go get -u $(DIR)/...

compile-server: deps
ifdef DEBUG
	@echo 'Compiling for debugging...'
	@$(GOARGS) go build $(GO_DEBUG_BUILD_ARGS) -o $(DIR)$(SERVER_BINARY_NAME) $(SERVER_BUILD_PATH)
else ifdef CONTAINER
	@echo 'Compiling for container usage...'
	@$(GOARGS) go build $(GO_CONTAINER_BUILD_ARGS) -o $(DIR)$(SERVER_BINARY_NAME) $(SERVER_BUILD_PATH)
else
	@echo 'Compiling for normal Linux env...'
	@$(GOARGS) go build $(GO_BUILD_ARGS) -o $(DIR)$(SERVER_BINARY_NAME) $(SERVER_BUILD_PATH)
endif

compile-cli: deps
	@$(GOARGS) go build $(GO_BUILD_ARGS) -o $(CLI_BINARY_NAME) $(CLI_BUILD_PATH)

compile: compile-server compile-cli

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
