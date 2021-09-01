# INetMock

[![pipeline status](https://gitlab.com/inetmock/inetmock/badges/main/pipeline.svg)](https://gitlab.com/inetmock/inetmock/-/commits/main)
[![coverage report](https://gitlab.com/inetmock/inetmock/badges/main/coverage.svg)](https://gitlab.com/inetmock/inetmock/-/commits/main)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/inetmock/inetmock)](https://goreportcard.com/report/gitlab.com/inetmock/inetmock)

INetMock is kind of a fork of [INetSim](https://www.inetsim.org/).
"Kind of" in terms that both applications overlap in their functionality to serve as "fake internet" routers.

INetMock right now does **not** implement so many protocols like INetSim. In fact it is only able to respond to HTTP,
HTTPS, DNS, DNS-over-TLS (DoT) requests and to act as an HTTP proxy. The most notable advantage of INetMock over INetSim
is that it issues proper TLS certificates on the fly signed by a CA certificate that can be deployed to client systems
to achieve "proper" TLS encryption - as long as the client does not use certificate pinning or something similar.

A second advantage is that INetMock is a complete rewrite in Go. It has a way smaller memory footprint and far better
startup and shutdown times. It also does not enforce `root` privileges as it is also possible to run the application
with the required capabilities to open ports e.g. with SystemD (a sample unit file can be found in the `deploy/`
directory).

_This project is still heavy work-in-progress. There may be breaking changes at any time. There's no guarantee for
anything except no kittens will be harmed!_

## Use cases

While the original use case was to simulate an internet connection both server and client might be used for other things too:

- serving as a mock API while developing an HTTP client library where you exactly know which requests should return which responses because you can match requests exactly with path and headers and return inline JSON, JSON from files, set status codes, ...
- serving as an advanced client CLI if you design an HTTP server application because you can run integration tests very easy including validation of results
- serving as an advanced client CLI if you design a custom DNS server because it's very easy to run queries (also from scripts) including support for custom ports - DoT and DoH client support is planned soon

## Qickstart

So you're asking 'how do I get started to see what this thing can do for me?!' - then this is for you! 

### Docker/Podman

The probably easiest way to get started is to use the pre-built container image.
The current tags can be found in the [releases](https://gitlab.com/inetmock/inetmock/-/releases).
The pre-built container image is configured with the [config-container.yaml](config-container.yaml) but you can always mount your own config.
Because the default config binds to the ports 53, 80 and 443 it requires some additional capabilities:

```
docker/podman run --rm -ti --cap-add CAP_NET_RAW --cap-add CAP_NET_BIND_SERVICE registry.gitlab.com/inetmock/inetmock:latest 
```

Depending on your use case it makes either sense to publish the ports of the container, run it in network mode `host` or isolate it to an internal network with the workload you're analyzing.
A very basic example how to run a Vagrant VM with an INetMock instanced running with Podman in a 'private' network can be found [here](https://gitlab.com/inetmock/examples/-/tree/master/vagrant-libvirt).

To run the container with a custom config just override the existing one like so:

```
docker/podman run --rm -ti -v `pwd`/config.yaml:/etc/inetmock/config.yaml:ro --cap-add CAP_NET_RAW --cap-add CAP_NET_BIND_SERVICE registry.gitlab.com/inetmock/inetmock:latest 
```

_Note:_ The pre-built container image is based on the 'distroless/static:nonroot'.
In consequence every file or directory you expect the container to access/modify needs corresponding access rights or you have to run the container with a different user.

### Binaries

Binaries can also be found on the [releases](https://gitlab.com/inetmock/inetmock/-/release) page.
Due to dependencies to some Linux sub-systems (e.g. the whole PCAP recording stuff) there is only a Linux binary of the INetMock server.
The client CLI `imctl` is available for Linux, MacOS and Windows (while it has to be noted that Windows and MacOS are not tested).

By default the server looks for `config.yaml` files in the following places:

- `/etc/inetmock/config.yaml`
- `$HOME/.inetmock/config.yaml`
- `./config.yaml`

Because INetMock requires a lot of setup it's not possible to configure it completely from flags hence you need a config in any of the aforementioned places.
If you don't know where to start the default `config.yaml` from this repository might be a good start because it's also the one that is used during development and therefore always up-to-date.

### `imctl`

To interact with the gRPC API of INetMock without having to write your own application `imctl` helps you to control your INetMock instance.
`imctl` can be used to (probably not exhaustive):

- interact with the audit API - the audit API allows you to monitor which requests INetMock handled in near-realtime, register an audit monitoring file to get a structured log, read those protobuf files to JSON and to remove an audit sink
- interact with the health API - runs the defined health checks on the server side and returns the result including an exit code != 0 if any check fails
- interact with the PCAP API - start/stop recording of network interface traffic to PCAP files, list available interfaces, list active recordings
- run check scripts like the [interation test](testdata/integration.imcs) or run single check commands like `imctl check run "http.GET('https://google.com/favicon.ico') => Status(200)"`

Everything that can be done from the CLI is documented with `--help` switches hence no huge documentation that will be outdated as soon as it is pushed here.

In general it always is a good idea to check the [Taskfile.yml](Taskfile.yml) and the [.gitlab-ci.yml](.gitlab-ci.yml) files for examples on how to use client and server for different use cases.

## Docs

Docs are available either in the [`docs/`](./docs/) directory or as rendered markdown book at
the [GitLab pages](https://inetmock.gitlab.io/inetmock/).

## Contribution/feature requests

Please create an issue for any proposal, feature requests, found bug,... I'm glad for every kind of feedback!

Right now I've no special workflow for pull requests but I will look into every proposed change.
