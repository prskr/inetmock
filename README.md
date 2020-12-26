# INetMock

![Go](https://gitlab.com/inetmock/inetmock/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/inetmock/inetmock)](https://goreportcard.com/report/gitlab.com/inetmock/inetmock)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=baez90_inetmock&metric=alert_status)](https://sonarcloud.io/dashboard?id=baez90_inetmock)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=baez90_inetmock&metric=ncloc)](https://sonarcloud.io/dashboard?id=baez90_inetmock)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=baez90_inetmock&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=baez90_inetmock)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=baez90_inetmock&metric=security_rating)](https://sonarcloud.io/dashboard?id=baez90_inetmock)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=baez90_inetmock&metric=bugs)](https://sonarcloud.io/dashboard?id=baez90_inetmock)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=baez90_inetmock&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=baez90_inetmock)

INetMock is kind of a fork of [INetSim](https://www.inetsim.org/).
"Kind of" in terms that both applications overlap in their functionality to serve as "fake internet" routers.

INetMock right now does **not** implement so many protocols like INetSim. In fact it is only able to respond to HTTP,
HTTPS, DNS, DNS-over-TLS (DoT) requests and to act as an HTTP proxy. The most notable advantage of INetMOck over INetSim
is that it issues proper TLS certificates on the fly signed by a CA certificate that can be deployed to client systems
to achieve "proper" TLS encryption - as long as the client does not use certificate pinning or something similar.

A second advantage is that INetMock is a complete rewrite in Go, based on a plugin system that allows dynamic
configuration while it has a way smaller memory footprint and far better startup and shutdown times. It also does not
enforce `root` privileges as it is also possible to run the application with the required capabilities to open ports
e.g. with SystemD (a sample unit file can be found in the `deploy/` directory).

_This project is still heavy work-in-progress. There may be breaking changes at any time. There's no guarantee for
anything except no kittens will be harmed!_

## Docs

Docs are available either in the [`docs/`](./docs/) directory or as rendered markdown book at
the [GitHub pages](https://baez90.github.io/inetmock/).

## Contribution/feature requests

Please create an issue for any proposal, feature requests, found bug,... I'm glad for every kind of feedback!

Right now I've no special workflow for pull requests but I will look into every proposed change.