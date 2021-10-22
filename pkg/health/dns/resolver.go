package dns

import (
	"context"
	"net"
)

type (
	Resolver interface {
		// LookupA looks up the given host using the corresponding query protocol.
		// It returns a slice of that host's addresses.
		LookupA(ctx context.Context, host string) (addrs []net.IP, err error)

		// LookupPTR performs a reverse lookup for the given address, returning a list
		// of names mapping to that address.
		//
		// The returned names are validated to be properly formatted presentation-format
		// domain names. If the response contains invalid names, those records are filtered
		// out and an error will be returned alongside the the remaining results, if any.
		LookupPTR(ctx context.Context, addr string) (names []string, err error)
	}
)
