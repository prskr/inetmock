package test

import (
	"gitlab.com/inetmock/inetmock/protocols/dns/client"
)

func DNSResolverForInMemListener(lis InMemListener) *client.Resolver {
	return &client.Resolver{
		Transport: &client.TraditionalTransport{
			Dial: lis.DialContext,
		},
	}
}
