package mock

import (
	"net"
	"regexp"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

type resolverRule struct {
	pattern  *regexp.Regexp
	response net.IP
}

type dnsOptions struct {
	Rules    []resolverRule
	Fallback ResolverFallback
}

func loadFromConfig(lifecycle endpoint.Lifecycle) (dnsOptions, error) {
	type rule struct {
		Pattern  string
		Response string
	}

	type fallback struct {
		Strategy string
		Args     map[string]interface{}
	}

	opts := struct {
		Rules    []rule
		Fallback fallback
	}{}

	if err := lifecycle.UnmarshalOptions(&opts); err != nil {
		return dnsOptions{}, err
	}

	var options dnsOptions
	for _, rule := range opts.Rules {
		var err error
		var rr resolverRule
		if rr.pattern, err = regexp.Compile(rule.Pattern); err != nil {
			continue
		}

		if rr.response = net.ParseIP(rule.Response); rr.response == nil {
			continue
		}
		options.Rules = append(options.Rules, rr)
	}

	options.Fallback = CreateResolverFallback(
		opts.Fallback.Strategy,
		opts.Fallback.Args,
	)

	return options, nil
}
