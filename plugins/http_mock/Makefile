VERSION = $(shell git describe --dirty --tags --always)
DIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
PKGS = $(shell go list ./...)
TEST_PKGS = $(shell go list ./...)
GOARGS = GOOS=linux GOARCH=amd64 CGO_ENABLED=0
GO_BUILD_ARGS = -buildmode=plugin -ldflags="-w -s"
GO_CONTAINER_BUILD_ARGS = -buildmode=plugin -ldflags="-w -s" -a -installsuffix cgo
GO_DEBUG_BUILD_ARGS = -buildmode=plugin -gcflags "all=-N -l"
PLUGIN_NAME = $(shell basename $(DIR)).so
OUT_DIR = $(DIR)../../
DEBUG_PORT = 2345

.PHONY: deps update-deps format compile test cli-cover-report html-cover-report

all: format compile test

deps:
	@go build -buildmode=plugin -v $(DIR)...

update-deps:
	@go mod tidy
	@go get -u

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
	@rm -f $(DIR)$(PLUGIN_NAME)

cli-cover-report:
	@go tool cover -func=cov.out

html-cover-report:
	@go tool cover -html=cov.out -o .coverage.html