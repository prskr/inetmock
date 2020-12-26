# `dns_mock`

## Intro

The `dns_mock` handler expects an array of rules how it should respond to dfferent DNS queries and a fallback strategy.
Currently only queries for __A__ records are supported. The rules are primarily meant to define some exceptions or well
known DNS responses e.g. to return to right Google DNS IP but for everything else it will return dummy IPs.

The rules for the `dns_mock` handler are equivalent to the `http_mock` rules. Every rule consists of a `pattern` that
specifies a query name e.g. a single host, a wildcard domain, a wildcard top-level domain or even a _"match all"_ rule
is possible. These rules are evaluated in the same order they are defined in the `config.yaml`.

The fallback strategy is taken into account whenever a query does not match a rule.

Right now the following fallback strategies are available:

* _random_
* _incremental_

Just like the handler is configured via the `options` object the fallback strategies are configured via an `args`
object.

### _random_ fallback

The _random_ fallback strategy is the easier one of the both. It doesn't take any argument and it just shuffles a random
IP address for every request no matter if it was already asked for this IP or not.

### _incremental_ fallback

The _incremental_ fallback is little bit more intelligent. It takes a `startIP` as an argument which defines from which
IP address the strategy starts counting up to respond to DNS queries. Just like the _incremental_ strategy it is _
stateless_ and does not store any already given response for later reuse (at least for now).

## Configuration

### Matching an explicit host

The easiest possible pattern is to match a single host:

```yml
endpoints:
  plainDns:
    handler: dns_mock
    listenAddress: 0.0.0.0
    port: 53
    options:
      rules:
        - pattern: "github\\.com"
          response: 1.1.1.1
```

### Matching a whole domain

While matching a single host is nice2have it's not very helpful in most cases - except for some edge cases where it
might be necesary to specifically return a certain IP address. But it's also possible to match a whole domain no matter
what subdomain or sub-subdomain or whatever is requested like this:

```yml
endpoints:
  plainDns:
    handler: dns_mock
    listenAddress: 0.0.0.0
    port: 53
    options:
      rules:
        - pattern: ".*\\.google\\.com"
          response: 2.2.2.2
```

### Matching a whole TLD

In some cases it might also be interesting to distinguish between different requested TLDs. Therefore it might be
interesting to define one IP address to resolve to for every TLD that should be distinguishable.

```yml
endpoints:
  plainDns:
    handler: dns_mock
    listenAddress: 0.0.0.0
    port: 53
    options:
      rules:
        - pattern: ".*\\.com"
          response: 2.2.2.2
```

### Matching any query

Last but not least it is obvously also possible to match any query. This is comparable to a _"static"_ fallback strategy
in cases where different IP addresses are not necessary but the network setup should be as easy as possible.

```yml
endpoints:
  plainDns:
    handler: dns_mock
    listenAddress: 0.0.0.0
    port: 53
    options:
      rules:
        - pattern: ".*"
          response: 10.0.10.1
```

### Fallback strategies

#### _random_

Like previously mentioned the _random_ strategy is easy as it can be. It just takes a random unsigned integer of 32
bits, converts it to an IP address and returns this address as response. Therefore no further configuration is necessary
for now.

```yml
endpoints:
  plainDns:
    handler: dns_mock
    listenAddress: 0.0.0.0
    port: 53
    options:
      rules: []
      fallback:
        strategy: random
```

#### _incremental_

Also like previously mentioned the _incremental_ fallback strategy is fairly easy to setup. It just takes a `startIP` as
argument which is used to count upwards. It does __not__ check for an interval or something like this right now so a
overflow might occur.

```yml
endpoints:
  plainDns:
    handler: dns_mock
    listenAddress: 0.0.0.0
    port: 53
    options:
      rules: []
      fallback:
        strategy: incremental
        args:
          startIP: 10.0.0.0
```