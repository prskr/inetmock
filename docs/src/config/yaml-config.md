# `config.yaml`

## Intro

The configuration of _INetMock_ is done in the `config.yaml`.
It defines which listeners should be started with which endpoint handler(s) and a few more things.

Every listener has a name that is used for logging and as already mentioned consists of protocol, port, optional an IP and the handler(s) and their options.

INetMock comes with _"Batteries included"_ and ships with a basic `config.yaml` that defines a basic set of endpoints for:

* HTTP
* HTTPS
* HTTP(S) proxy
* DNS
* DNS-over-TLS
* Prometheus metrics

Because the config is YAML it's also possible to use YAML anchors to deduplicate the configuration like shown in the included default configuration.

## Multiplexing

Endpoint multiplexing is still a very new feature.
Currently only the `http_mock` and the `http_proxy` handler support multiplexing but because both are HTTP handlers they can only be combined in special circumstances.