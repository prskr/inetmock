package mock

import (
	"net"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
)

type cache interface {
	PutRecord(host string, address net.IP)
	ForwardLookup(host string, resolver dns.IPResolver) net.IP
	ReverseLookup(address net.IP) (host string, miss bool)
}
