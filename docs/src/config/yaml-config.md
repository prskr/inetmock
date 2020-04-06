
# `config.yaml`

## Intro

The configuration of _INetMock_ is mostly done in the `config.yaml`.
It defines which endpoints should be started with which handler and a few more things.

Every endpoint has a name that is used for logging and as already mentioned consists of listening IP and port, the handler and its options.

INetMock comes with _"Batteries included"_ and ships with a basic `config.yaml` that defines a basic set of endpoints for:

* HTTP
* HTTPS
* DNS
* DNS-over-TLS

## Sample

```yml
endpoints:
    myHttpEndpoint:
        handler: http_mock
        listenAddress: 127.0.0.1
        port: 8080
        options:
            rules:
                - pattern: ".*"
                  target: ./assets/fakeFiles/default.html
```