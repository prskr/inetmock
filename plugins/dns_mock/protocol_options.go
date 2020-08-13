package dns_mock

import (
	"github.com/spf13/viper"
	"net"
	"regexp"
)

const (
	fallbackArgsConfigKey = "fallback.args"
)

type resolverRule struct {
	pattern  *regexp.Regexp
	response net.IP
}

type dnsOptions struct {
	Rules    []resolverRule
	Fallback ResolverFallback
}

func loadFromConfig(config *viper.Viper) (options dnsOptions) {
	type rule struct {
		Pattern  string
		Response string
	}

	type fallback struct {
		Strategy string
	}

	opts := struct {
		Rules    []rule
		Fallback fallback
	}{}

	_ = config.Unmarshal(&opts)

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
		config.Sub(fallbackArgsConfigKey),
	)

	return
}
