package mock

import (
	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
)

type dnsOptions struct {
	Rules   []string
	Cache   cache
	Default dns.IPResolver
}

func loadFromConfig(lifecycle endpoint.Lifecycle) (dnsOptions, error) {
	var (
		opts dnsOptions
	)
	if err := lifecycle.UnmarshalOptions(&opts); err != nil {
		return dnsOptions{}, err
	}

	return opts, nil
}
