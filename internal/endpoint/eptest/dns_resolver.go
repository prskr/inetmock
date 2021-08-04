package eptest

import (
	"net"
)

func DNSResolverForInMemListener(lis InMemListener) *net.Resolver {
	return &net.Resolver{
		PreferGo:     true,
		Dial:         lis.DialContext,
		StrictErrors: true,
	}
}
