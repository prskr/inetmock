package mock

import (
	"net"
	"sync/atomic"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
)

type IncrementalIPResolver struct {
	offset uint32
	cidr   *net.IPNet
}

func (i *IncrementalIPResolver) Lookup(string) net.IP {
	var (
		ones, bits   = i.cidr.Mask.Size()
		max          = uint32(1<<(bits-ones)) - 1
		offset, base uint32
	)

	atomic.CompareAndSwapUint32(&i.offset, max, 0)
	offset = atomic.AddUint32(&i.offset, 1)

	base = dns.IPToInt32(i.cidr.IP)
	return dns.Uint32ToIP(base + offset)
}
