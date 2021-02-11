# Configuration

## Listeners

_INetMock_ runs _listeners_ that consist basically of a _port_ and a _protocol_ and have at least one _endpoint_ that handles incoming requests.

The basic configuration looks like this:

```yaml
listeners:
    [name]:
        protocol: tcp
        port: 80
        endpoints: {}
```

optionally the listener can be restricted to an IP address:

```yaml
listeners:
    [name]:
        listenAddress: 127.0.0.1
        protocol: tcp
        port: 80
        endpoints: {}
    [name]:
        listenAddress: [::1]
        protocol: tcp
        port: 80
        endpoints: {}
```

Supported protocols are

* `tcp`
* `udp`

_Note:_ technically Go would also support something like `tcp4`, `tcp6` and a few more but in most cases they shouldn't be necessary.

## Endpoints

_INetMock_ ships multiple endpoint handlers:

* [HTTP mock `http_mock`](config/http_mock.md)
* [HTTP proxy `http_proxy`](config/http_proxy.md)
* [DNS mock `dns_mock`](config/dns_mock.md)
* [Metrics exporter `metrics_exporter`](config/metrics_exporter.md)

Some of the endpoints can be combined e.g. it's possible to run two instances of the `http_mock` handler on the same _listener_ one with TLS enabled and one without to support HTTP and HTTPS at the same time on the same port.

Details can be found in the doc pages of every handler.

## TLS

One of the more interesting aspects of _INetMock_ is that it generates TLS certificates on-demand and (more interesting) **on-the-fly**.

This allows [_man-in-the-middle (MITM)_ attacks](https://en.wikipedia.org/wiki/Man-in-the-middle_attack) as long as the client trusts the root certificate used by _INetMock_.

_INetMock_ ships with a default CA certificate in the container image that can be found in the `assets/demoCA` directory in the [repository](https://gitlab.com/inetmock/inetmock).
The default CA certificate is valid from the year _02-01-2001_ to _01-20-2051_ to allow also the creation of certificates in the future or the past (see [features](features.md) for planned features).

Certificates are cached and re-used if applicable as long as they are still valid.
See the [TLS config section](config/tls.md) for further details.

## API

_INetMock_ comes with a gRPC API that by default listens only on a Unix socket in the path `/var/run/inetmock.sock`

The API can be used either programmatically with any gRPC enabled language based on the `.proto` files in the `api/proto` directory of the [repository](https://gitlab.com/inetmock/inetmock) or with the `imctl` CLI.

Right now the API possibilities are rather limited because it only offers options to use the more mock alike health API and to interact with the [Audit API](features/audit.md) but there's more to come.