# INetMock

[![pipeline status](https://gitlab.com/inetmock/inetmock/badges/master/pipeline.svg)](https://gitlab.com/inetmock/inetmock/-/commits/master)
[![coverage report](https://gitlab.com/inetmock/inetmock/badges/master/coverage.svg)](https://gitlab.com/inetmock/inetmock/-/commits/master)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/inetmock/inetmock)](https://goreportcard.com/report/gitlab.com/inetmock/inetmock)

INetMock is kind of a fork of [INetSim](https://www.inetsim.org/).
"Kind of" in terms that both applications overlap in their functionality to serve as "fake internet" routers.

INetMock right now does **not** implement so many protocols like INetSim. In fact it is only able to respond to HTTP,
HTTPS, DNS, DNS-over-TLS (DoT) requests and to act as an HTTP proxy. The most notable advantage of INetMOck over INetSim
is that it issues proper TLS certificates on the fly signed by a CA certificate that can be deployed to client systems
to achieve "proper" TLS encryption - as long as the client does not use certificate pinning or something similar.

A second advantage is that INetMock is a complete rewrite in Go. It has a way smaller memory footprint and far better
startup and shutdown times. It also does not enforce `root` privileges as it is also possible to run the application
with the required capabilities to open ports e.g. with SystemD (a sample unit file can be found in the `deploy/`
directory).

_This project is still heavy work-in-progress. There may be breaking changes at any time. There's no guarantee for
anything except no kittens will be harmed!_

## Docs

Docs are available either in the [`docs/`](./docs/) directory or as rendered markdown book at
the [GitHub pages](https://baez90.github.io/inetmock/).

## Contribution/feature requests

Please create an issue for any proposal, feature requests, found bug,... I'm glad for every kind of feedback!

Right now I've no special workflow for pull requests but I will look into every proposed change.