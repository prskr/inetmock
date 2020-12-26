# Configuration

## Plugins & handlers

_INetMock_ is based on plugins that ship one or more __protocol handlers__. Examples for protocol handlers are HTTP or
DNS but also TLS.

The application ships with the following handlers:

* `http_mock`
* `dns_mock`
* `tls_interceptor`

The configuration of an so called endpoint always specifies which handler should be used, which IP address and port it
should listen on and some handler specific `options`. This way the whole system is very flexible and can be configured
for various individual scenarios.

## Commands

Beside of __protocol handlers__ a plugin can also ship custom commands e.g. the `tls_interceptor` ships a `generate-ca`
command to bootstrap a certificate authority key-pair that can be reused for multiple instances.