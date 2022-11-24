package test

import (
	"inetmock.icb4dc0.de/inetmock/protocols/dns/client"
)

func DNSResolverForInMemListener(lis InMemListener) *client.Resolver {
	return &client.Resolver{
		Transport: &client.TraditionalTransport{
			Dial: lis.DialContext,
		},
	}
}
