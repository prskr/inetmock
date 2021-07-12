package mock

import (
	"math/rand"
	"net"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
)

type RandomIPResolver struct {
	random *rand.Rand
}

func (r *RandomIPResolver) Lookup(string) net.IP {
	return dns.Uint32ToIP(r.random.Uint32())
}
