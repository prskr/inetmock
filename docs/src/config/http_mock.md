# `http_mock`

## Intro

The `http_mock` handler expects an array of rules how it should respond to different request paths.
This allows to e.g. return an image if the request path contains something like _"asdf.jpg"_ but with binary if the request path contains something like _"malicous.exe"_.

A _"catch all"_ rule could return in any case an HTML page or if nothing is provided the handler returns an HTTP 404 status code.

The rules are taken into account in the same order than they are defined in the `config.yaml`.

Every rule consists of a regex `pattern` (__re2__ compatible) and a `response` path to the file it should return.

In the future more advanced rules might be possible e.g. to match not on the request path but on some header values.

## Configuration

### Matching a specific path

The easiest possible pattern is to match a static request path:

```yml
endpoints:
  plainHttp:
    handler: http_mock
    listenAddress: 0.0.0.0
    port: 80
    options:
      rules:
        - pattern: "/static/http/path/sample.exe"
          response: ./assets/fakeFiles/sample.exe
```

### Matching a file extensions

While matching a static path might be nice as an example it's not very useful.
Returning a given file for all kinds of of request paths based on the requested file extension is way more handy:

```yml
endpoints:
  plainHttp:
    handler: http_mock
    listenAddress: 0.0.0.0
    port: 80
    options:
      rules:
        - pattern: ".*\\.png"
          response: ./assets/fakeFiles/default.png
```

So this is already way more flexible but we can do even better:

```yml
endpoints:
  plainHttp:
    handler: http_mock
    listenAddress: 0.0.0.0
    port: 80
    options:
      rules:
        - pattern: ".*\\.(?i)(jpg|jpeg)"
          response: ./assets/fakeFiles/default.jpg
```

This way the extension ignores any case and matches both `.jpg` and `.jpeg` (and of course also e.g. `.JpEg` and so on and so forth).

The default `config.yaml` already ships with some basic rules to handle the most common file extensions.

### Defining a fallback

Last but not least a default case might be necessary to get at least any response but a 404.

This can be achieved with a `.*` pattern that literally matches everything:

```yml
endpoints:
  plainHttp:
    handler: http_mock
    listenAddress: 0.0.0.0
    port: 80
    options:
      rules:
        - pattern: ".*"
          response: ./assets/fakeFiles/default.html
```