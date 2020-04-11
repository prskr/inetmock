VERSst/pluginsION = $(shell git describe --dirty --tags --always)
DIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
BUILD_PATH = $(DIR)/main.go
PKGS = $(shell go list ./...)
TEST_PKGS = $(shell find . -type f -name "*_test.go" -not -path "./plugins/*" -printf '%h\n' | sort -u)
GOARGS = GOOS=linux GOARCH=amd64
GO_BUILD_ARGS = -ldflags="-w -s"
GO_CONTAINER_BUILD_ARGS = -ldflags="-w -s" -a -installsuffix cgo
GO_DEBUG_BUILD_ARGS = -gcflags "all=-N -l"
BINARY_NAME = inetmock
PLUGINS = $(wildcard $(DIR)plugins/*/.)
DEBUG_PORT = 2345
DEBUG_ARGS?= --development-logs=true
INETMOCK_PLUGINS_DIRECTORY = $(DIR)

.PHONY: clean all format deps compile debug snapshot-release test cli-cover-report html-cover-report plugins $(PLUGINS)

all: clean format compile test plugins

clean:
	@find $(DIR) -type f \( -name "*.out" -or -name "*.so" \) -exec rm -f {} \;
	@rm -rf $(DIR)*.so
	@rm -f $(DIR)$(BINARY_NAME) $(DIR)main

format:
	@go fmt $(PKGS)

deps:
	@go mod tidy
	@go get -u
	@go build -v $(BUILD_PATH)

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

debug:
	@export INETMOCK_PLUGINS_DIRECTORY
	@dlv exec $(DIR)$(BINARY_NAME) \
		--headless \
		--listen=:2345 \
		--api-version=2 \
		--accept-multiclient \
		-- $(DEBUG_ARGS)

snapshot-release:
	@goreleaser release --snapshot --skip-publish --rm-dist

test:
	@go test -coverprofile=./cov-raw.out -v $(TEST_PKGS)
	@cat ./cov-raw.out | grep -v "generated" > ./cov.out

cli-cover-report:
	@go tool cover -func=cov.out

html-cover-report:
	@go tool cover -html=cov.out -o .coverage.html

plugins: $(PLUGINS)
$(PLUGINS):
	$(MAKE) -C $@