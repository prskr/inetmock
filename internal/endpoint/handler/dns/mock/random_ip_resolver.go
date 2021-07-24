package mock

import (
	"math/rand"
	"net"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
)

type RandomIPResolver struct {
	Random *rand.Rand
	CIDR   *net.IPNet
}

func (r *RandomIPResolver) Lookup(string) net.IP {
	var (
		ones, bits   = r.CIDR.Mask.Size()
		max          = (1 << (bits - ones)) - 1
		offset, base uint32
	)

	base = dns.IPToInt32(r.CIDR.IP)
	offset = uint32(r.Random.Intn(max))

	return dns.Uint32ToIP(base + offset)
}
