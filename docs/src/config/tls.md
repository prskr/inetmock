# TLS

The default TLS configuration looks like the following:

```yaml
tls:
  curve: P256
  minTLSVersion: SSL3
  includeInsecureCipherSuites: false
  validity:
    ca:
      notBeforeRelative: 17520h
      notAfterRelative: 17520h
    server:
      NotBeforeRelative: 168h
      NotAfterRelative: 168h
  rootCaCert:
    publicKeyPath: ./assets/demoCA/ca.pem
    privateKeyPath: ./assets/demoCA/ca.key
  certCachePath: /tmp/inetmock/
```

In the following sections every aspects will be explained in detail.

## `curve`

For different reasons _INetMock_ enforces _Elliptic-curve cryptography (ECC)_.
Besides the better security ECC certificates are way smaller, faster to generate and faster in the usage than [RSA](https://en.wikipedia.org/wiki/RSA_(cryptosystem)) certificates.

`curve` configures the ECC algorithm to use for the creation of ephemeral server certificates.
Possible values are:

* P224
* P256
* P384
* P521

If the curve is malformed configured _INetMock_ will fallback to **P256**.

## `minTLSVersion`

`minTLSVersion` configures the minimal allowed TLS version for all TLS connections.
Possible values are:

* SSL3
* TLS10
* TLS11
* TLS12
* TLS13

If either an invalid or an empty value are set _INetMock_ will fallback to **TLS13**.
This might change at any time.

In the default configuration above it is set to SSL3 to allow even old browsers/devices to connect as long as they support ECC based cipher suites.

## `includeInsecureCipherSuites`

`includeInsecureCipherSuites` enforces to include cipher suites that are known to be insecure to support even more old clients.

## `validity`

The `validity` section configures the validity range of the ephemeral server certificates (and it allows to set the validity of the CA cert when it's generated but that's covered in the [CLIs section](../features/cli.md)).

The `validity` consists of a `NotBeforeRelative` and a `NotAfterRelative`.
The actual start and end times are calculated right when the certificate is generated.
This allows to create certificates with a validity in the past or the future depending on what _INetMock_ takes as the current time.

As already mentioned in the general config section certificates are cached and recreated if they are invalid or "about to become invalid".

"About to become invalid" is calculated as relative time span based on the current time and the total validity range.
If the certificate is not longer valid than 25% of its total validity time it's considered to be invalid and will be recreated.

An example might be helpful:

If `NotBeforeRelative` and `NotAfterRelative` are both `12h` the certificate will be considered invalid if it's remaining validity time is less than `6h`.

### Units

All time units are parsed as Go's `time.Duration`.
See the [docs](https://golang.org/pkg/time/#ParseDuration) for more details which units are supported.

## `rootCaCert`

`rootCaCert` configures the paths to the public and private key for the root CA certificate used to sign the ephemeral server certificates.

## `certCachePath`

`certCachePath` configures the path where the ephemeral certificates are kept for reusing.