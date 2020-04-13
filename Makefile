VERSst/pluginsION = $(shell git describe --dirty --tags --always)
DIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
BUILD_PATH = $(DIR)/main.go
PKGS = $(shell go list ./...)
TEST_PKGS = $(shell find . -type f -name "*_test.go" -not -path "./plugins/*" -printf '%h\n' | sort -u)
GO_GEN_FILES = $(shell grep -rnwl --include="*.go" "go:generate" $(DIR))
GOARGS = GOOS=linux GOARCH=amd64
GO_BUILD_ARGS = -ldflags="-w -s"
GO_CONTAINER_BUILD_ARGS = -ldflags="-w -s" -a -installsuffix cgo
GO_DEBUG_BUILD_ARGS = -gcflags "all=-N -l"
BINARY_NAME = inetmock
PLUGINS = $(wildcard $(DIR)plugins/*/.)
DEBUG_PORT = 2345
DEBUG_ARGS?= --development-logs=true

.PHONY: clean all format deps update-deps compile debug generate snapshot-release test cli-cover-report html-cover-report plugins $(PLUGINS) $(.PHONY: clean all format deps update-deps compile debug generate snapshot-release test cli-cover-report html-cover-report plugins $(PLUGINS) $(GO_GEN_FILES)
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
	dlv exec $(DIR)$(BINARY_NAME) \
		--headless \
		--listen=:2345 \
		--api-version=2 \
		--accept-multiclient \
		-- $(DEBUG_ARGS)

generate:
	@for go_gen_target in $(GO_GEN_FILES); do \
  		go generate $$go_gen_target; \
  	done

snapshot-release:
	@goreleaser release --snapshot --skip-publish --rm-dist

test:
	@go test -coverprofile=./cov-raw.out -v $(TEST_PKGS)
	@cat ./cov-raw.out | grep -v "generated" > ./cov.out

cli-cover-report:
	@go tool cover -func=cov.out

html-cover-report: test
	@go tool cover -html=cov.out -o .coverage.html

plugins: $(PLUGINS)
$(PLUGINS):
	$(MAKE) -C $@