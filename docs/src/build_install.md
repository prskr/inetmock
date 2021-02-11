# Installation

The recommended way to get INetMock is to download either a [release](https://gitlab.com/inetmock/inetmock/-/releases) or to use it as [container](https://gitlab.com/inetmock/inetmock/container_registry/1605679).

## Building from source

### Requirements

Mandatory:

* Go 1.15
* [mockgen](https://github.com/golang/mock/)
* [go-enum](https://github.com/abice/go-enum)
* protoc
* [protoc-gen-go](https://developers.google.com/protocol-buffers/docs/reference/go-generated)
* [protoc-gen-go-grpc](https://grpc.io/docs/languages/go/quickstart/)

Optional/Development:

* [GoReleaser](https://goreleaser.com/)
* [Task](https://taskfile.dev/#/)
* [golangci-lint](https://github.com/golangci/golangci-lint)

### Binaries

INetMock consists of two components:

* the _inetmock_ server that is the actual server responding to requests
* the _imctl_ client CLI to interact with the server e.g. to attach to audit logs

Both components can be downloaded as pre-built binaries from the [releases](https://gitlab.com/inetmock/inetmock/-/releases) which is the **recommended** way to get them.

#### _inetmock_ server

The easiest way to build the server is either to run `task build-inetmock` or `task snapshot-release`.
In both cases the `task` will take care that all code generators and the test suite will be executed at first to make sure that everything is fine.

If you don't want to use `task` or you can't for some reason the following steps have to be executed:

```sh
go generate ./...
protoc \
    --proto_path ./api/proto/ \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    $(find ./api/ -type f -name "*.proto" -printf "%p ")
go build -o inetmock ./cmd/inetmock
```

The above script creates a **non-optimized** binary to run the server.

#### _imctl_ client

The easiest way to build the server is either to run `task build-imctl` or `task snapshot-release`.
In both cases the `task` will take care that all code generators and the test suite will be executed at first to make sure that everything is fine.

If you don't want to use `task` or you can't for some reason the following steps have to be executed:

```sh
go generate ./...
protoc \
    --proto_path ./api/proto/ \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    $(find ./api/ -type f -name "*.proto" -printf "%p ")
go build -o imctl ./cmd/imctl
```

The above script creates a **non-optimized** binary to run the CLI.

## Docker/Podman

Container images can be found in the [GitLab container registry](https://gitlab.com/inetmock/inetmock/container_registry/1605679).